export function formatTimestamp(timestamp: Date | string): string {
    if (!timestamp) {
        return '';
    }
    
    const date = timestamp instanceof Date ? timestamp : new Date(timestamp);
    
    if (isNaN(date.getTime())) {
        return '';
    }
    
    return date.toLocaleTimeString('en-US', {
        hour: 'numeric',
        minute: '2-digit',
        hour12: true
    });
}

export function formatChatDate(timestamp: Date | string | null): string {
    if (!timestamp) {
        return '';
    }
    
    const date = timestamp instanceof Date ? timestamp : new Date(timestamp);
    
    if (isNaN(date.getTime())) {
        return '';
    }
    
    const now = new Date();
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const yesterday = new Date(today);
    yesterday.setDate(yesterday.getDate() - 1);
    const messageDate = new Date(date.getFullYear(), date.getMonth(), date.getDate());
    
    if (messageDate.getTime() === today.getTime()) {
        // Today - show time
        return date.toLocaleTimeString('en-US', {
            hour: 'numeric',
            minute: '2-digit',
            hour12: true
        });
    } else if (messageDate.getTime() === yesterday.getTime()) {
        // Yesterday
        return 'Yesterday';
    } else {
        // Older - show date
        const daysDiff = Math.floor((today.getTime() - messageDate.getTime()) / (1000 * 60 * 60 * 24));
        if (daysDiff < 7) {
            // Within a week - show day name
            return date.toLocaleDateString('en-US', { weekday: 'short' });
        } else {
            // Older - show date
            return date.toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric'
            });
        }
    }
}
