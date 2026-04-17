<script lang="ts">
    import type { Message } from '$lib/models/message';
    import { formatTimestamp, formatTimestampFull } from '$lib/utils/date';
    import checkIcon from '$lib/assets/icons/check.svg?url';
    import doubleCheckIcon from '$lib/assets/icons/check_double.svg?url';

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
            <div
                class="message-timestamp"
                title={formatTimestampFull(message.createdAt)}
                aria-label="Sent at {formatTimestamp(message.createdAt)}"
            >
                {formatTimestamp(message.createdAt)}
            </div>
            {#if message.fromMe}
                <span class="message-status" aria-label={message.readAt ? 'Read' : 'Sent'} title={message.readAt ? 'Read' : 'Sent'}>
                    <img
                        src={message.readAt ? doubleCheckIcon : checkIcon}
                        alt=""
                        aria-hidden="true"
                    />
                </span>
            {/if}
        </div>
    </div>
</div>
