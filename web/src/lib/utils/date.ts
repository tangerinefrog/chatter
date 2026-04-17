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

export function formatDateShort(timestamp: Date | string): string {
    if (!timestamp) {
        return '';
    }

    const date = timestamp instanceof Date ? timestamp : new Date(timestamp);
    if (isNaN(date.getTime())) {
        return '';
    }

    return date.toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric'
    });
}

export function formatTimestampFull(timestamp: Date | string): string {
    if (!timestamp) {
        return '';
    }

    const date = timestamp instanceof Date ? timestamp : new Date(timestamp);

    if (isNaN(date.getTime())) {
        return '';
    }

    return date.toLocaleString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: 'numeric',
        minute: '2-digit',
        second: '2-digit',
        hour12: true
    });
}

export function isSameDay(a: Date, b: Date): boolean {
    return (
        a.getFullYear() === b.getFullYear() &&
        a.getMonth() === b.getMonth() &&
        a.getDate() === b.getDate()
    );
}

export function isToday(timestamp: Date | string): boolean {
    if (!timestamp) {
        return false;
    }

    const date = timestamp instanceof Date ? timestamp : new Date(timestamp);
    if (isNaN(date.getTime())) {
        return false;
    }

    const today = new Date();
    return (
        date.getFullYear() === today.getFullYear() &&
        date.getMonth() === today.getMonth() &&
        date.getDate() === today.getDate()
    );
}
