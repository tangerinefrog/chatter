<script lang="ts">
    import type { Message } from '$lib/models/message';
    import { formatTimestamp, formatTimestampFull } from '$lib/utils/date';
    import { getFileUrl } from '$lib/api/client';
    import checkIcon from '$lib/assets/icons/check.svg?url';
    import doubleCheckIcon from '$lib/assets/icons/check_double.svg?url';
    import downloadIcon from '$lib/assets/icons/download.svg?url';
    import { isImage, formatFileSize, downloadFile } from '$lib/utils/files';
    import { onMount, onDestroy } from 'svelte';
    import PhotoSwipeLightbox from 'photoswipe/lightbox';
    import 'photoswipe/style.css';

    export let message: Message;
    export let chatId: string;

    let galleryId = `pswp-gallery-${message.id}`;
    let lightbox: PhotoSwipeLightbox | null = null;

    function setPreviewSize(event: Event) {
        const img = event.currentTarget as HTMLImageElement;
        const link = img.closest('a');

        if (link && img.naturalWidth && img.naturalHeight) {
            link.dataset.pswpWidth = String(img.naturalWidth);
            link.dataset.pswpHeight = String(img.naturalHeight);
        }
    }

    onMount(() => {
        if (!message.files?.some((file) => isImage(file.mimeType))) {
            return;
        }

        lightbox = new PhotoSwipeLightbox({
            gallery: `#${galleryId}`,
            children: 'a',
            pswpModule: () => import('photoswipe')
        });

        lightbox.init();
    });

    onDestroy(() => {
        lightbox?.destroy();
    });
</script>

<div
    class="message-row {message.fromMe ? 'me' : 'them'}"
    data-message-id={message.id}
    data-message-date={message.createdAt.toISOString()}
    data-read={message.readAt ? 'true' : 'false'}
>
    <div class="bubble">
        {#if message.text}
            <div class="message-text">{message.text}</div>
        {/if}

        {#if message.files && message.files.length > 0}
            <div class="files-container" id={galleryId}>
                {#each message.files as file}
                    {#if isImage(file.mimeType)}
                        <a
                            class="file-preview-link"
                            href={getFileUrl(chatId, file.id)}
                            data-pswp-width="1200"
                            data-pswp-height="900"
                            title={file.name}
                        >
                            <img
                                src={getFileUrl(chatId, file.id)}
                                alt={file.name}
                                class="file-preview"
                                title={file.name}
                                on:load={setPreviewSize}
                            />
                        </a>
                    {:else}
                        <a
                            class="file-item file-download-link"
                            href={getFileUrl(chatId, file.id)}
                            download={file.name}
                            aria-label="Download {file.name}"
                            title="Download {file.name}"
                        >
                            <div class="file-info">
                                <div class="file-name">{file.name}</div>
                                <div class="file-size">{formatFileSize(file.sizeBytes)}</div>
                            </div>
                            <img src={downloadIcon} alt="Download" aria-hidden="true" />
                        </a>
                    {/if}
                {/each}
            </div>
        {/if}

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
