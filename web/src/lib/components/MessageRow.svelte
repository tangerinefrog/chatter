<script lang="ts">
    import type { Message } from '$lib/models/message';
    import { formatTimestamp } from '$lib/utils/date';

    export let message: Message;
</script>

<div
    class="message-row {message.fromMe ? 'me' : 'them'}"
    data-message-id={message.id}
    data-read={message.readAt ? 'true' : 'false'}
>
    <div class="bubble">
        <div class="message-text">{message.text}</div>
        <div class="message-footer">
            <div class="message-timestamp" aria-label="Sent at {formatTimestamp(message.createdAt)}">
                {formatTimestamp(message.createdAt)}
            </div>
            {#if message.fromMe}
                <span class="message-status" aria-label="{message.readAt ? 'Read' : 'Sent'}">
                    {#if message.readAt}
                        ✓✓
                    {:else}
                        ✓
                    {/if}
                </span>
            {/if}
        </div>
    </div>
</div>
