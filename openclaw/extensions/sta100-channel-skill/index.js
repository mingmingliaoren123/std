import fs from "node:fs";
import path from "node:path";
import { request } from "node:http";

function readToken() {
  if (process.env.STA100_CHANNEL_SKILL_TOKEN?.trim()) return process.env.STA100_CHANNEL_SKILL_TOKEN.trim();
  const stateDir = process.env.OPENCLAW_STATE_DIR?.trim();
  if (!stateDir) return "";
  try {
    return fs.readFileSync(path.resolve(stateDir, "..", "channel-skill.token"), "utf8").trim();
  } catch {
    return "";
  }
}

function postInbound(payload, token, logger) {
  const endpoint = process.env.STA100_CHANNEL_SKILL_URL?.trim() || "http://127.0.0.1:18080/api/v1/openclaw/inbound";
  const url = new URL(`${endpoint.replace(/\/$/, "")}/${encodeURIComponent(payload.channel)}`);
  const body = JSON.stringify(payload);
  const req = request(url, {
    method: "POST",
    timeout: 3000,
    headers: {
      "content-type": "application/json",
      "content-length": Buffer.byteLength(body),
      "x-sta100-channel-token": token
    }
  }, (res) => {
    res.resume();
    if (res.statusCode && res.statusCode >= 400) logger.warn?.(`sta100-channel-skill: inbound handoff rejected (${res.statusCode})`);
  });
  req.on("timeout", () => req.destroy());
  req.on("error", (error) => logger.debug?.(`sta100-channel-skill: inbound handoff unavailable (${String(error)})`));
  req.end(body);
}

export default {
  id: "sta100-channel-skill",
  name: "STA-100 Channel Skill",
  description: "Forwards authorized channel messages to STA-100's confirmation workflow.",
  register(api) {
    api.logger.info?.("sta100-channel-skill: registering message_received hook");
    api.on("message_received", async (event, ctx) => {
      const token = readToken();
      if (!token || !ctx?.channelId) {
        api.logger.debug?.("sta100-channel-skill: skipped message without token or channel context");
        return;
      }
      const metadata = event.metadata && typeof event.metadata === "object" ? event.metadata : {};
      const content = event.content || event.body || event.bodyForAgent || "";
      postInbound({
        channel: ctx.channelId,
        account: ctx.accountId || event.accountId || "",
        conversation: ctx.conversationId || event.conversationId || metadata.conversationId || "",
        sender: event.senderId || event.from || "",
        from: event.from || "",
        messageId: event.messageId || "",
        message: content,
        timestamp: event.timestamp || 0,
        metadata: { threadId: event.threadId || metadata.threadId || "", sessionKey: event.sessionKey || ctx.sessionKey || "" }
      }, token, api.logger);
    }, { name: "sta100-channel-skill-message-received" });
    api.logger.info?.("sta100-channel-skill: message_received hook registered");
  }
};
