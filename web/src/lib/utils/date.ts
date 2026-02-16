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