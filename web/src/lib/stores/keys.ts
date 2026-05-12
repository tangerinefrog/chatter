import { deriveChatKey } from '$lib/crypto';

const STORAGE_KEY = 'chatter_master_key';

async function setMasterKey(key: Uint8Array): Promise<void> {
	const hex = Array.from(key).map(b => b.toString(16).padStart(2, '0')).join('');
	localStorage.setItem(STORAGE_KEY, hex);
}

function getMasterKey(): Uint8Array | null {
	const hex = localStorage.getItem(STORAGE_KEY);
	if (!hex) {
		return null;
	}
	return Uint8Array.from(hex.match(/.{1,2}/g)!.map(byte => parseInt(byte, 16)));
}

async function getChatKey(chatId: number): Promise<Uint8Array | null> {
	const masterKey = getMasterKey();
	if (!masterKey) {
		return null;
	}
	return await deriveChatKey(masterKey, chatId);
}

function clearKeys(): void {
	localStorage.removeItem(STORAGE_KEY);
}

export { setMasterKey, getMasterKey, getChatKey, clearKeys };
