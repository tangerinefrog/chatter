async function deriveChatKey(masterKey: Uint8Array, chatId: number): Promise<Uint8Array> {
	const encoder = new TextEncoder();
	const chatIdStr = chatId.toString();
	
	const combined = new Uint8Array(masterKey.length + chatIdStr.length);
	combined.set(masterKey, 0);
	combined.set(encoder.encode(chatIdStr), masterKey.length);
	
	const hashBuffer = await crypto.subtle.digest('SHA-256', combined);
	return new Uint8Array(hashBuffer);
}

async function encryptMessage(key: Uint8Array, plaintext: string): Promise<string> {
	const encoder = new TextEncoder();
	const plaintextData = encoder.encode(plaintext);
	
	const nonce = crypto.getRandomValues(new Uint8Array(12));
	
	const cryptoKey = await crypto.subtle.importKey(
		'raw',
		new Uint8Array(key) as BufferSource,
		{ name: 'AES-GCM', length: 256 },
		false,
		['encrypt']
	);
	
	const ciphertext = await crypto.subtle.encrypt(
		{ name: 'AES-GCM', iv: nonce },
		cryptoKey,
		plaintextData
	);
	
	const combined = new Uint8Array(nonce.length + ciphertext.byteLength);
	combined.set(nonce, 0);
	combined.set(new Uint8Array(ciphertext), nonce.length);
	
	return btoa(String.fromCharCode(...combined));
}

async function decryptMessage(key: Uint8Array, encoded: string): Promise<string> {
	const combined = Uint8Array.from(atob(encoded), c => c.charCodeAt(0));
	
	const nonce = combined.slice(0, 12);
	const ciphertext = combined.slice(12);
	
	const cryptoKey = await crypto.subtle.importKey(
		'raw',
		new Uint8Array(key) as BufferSource,
		{ name: 'AES-GCM', length: 256 },
		false,
		['decrypt']
	);
	
	const plaintext = await crypto.subtle.decrypt(
		{ name: 'AES-GCM', iv: nonce },
		cryptoKey,
		ciphertext
	);
	
	const decoder = new TextDecoder();
	return decoder.decode(plaintext);
}

function hexToUint8Array(hex: string): Uint8Array {
	const bytes = new Uint8Array(hex.length / 2);
	for (let i = 0; i < hex.length; i += 2) {
		bytes[i / 2] = parseInt(hex.substr(i, 2), 16);
	}
	return bytes;
}

export { deriveChatKey, encryptMessage, decryptMessage, hexToUint8Array };
