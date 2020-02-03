function str2buf(str) {
        return new TextEncoder("utf-8").encode(str);
}

function hex2buf(hexStr) {
        return new Uint8Array(hexStr.match(/.{2}/g).map(h => parseInt(h, 16)));
}

function buf2str(buffer) {
        return new TextDecoder("utf-8").decode(buffer);
}

async function deriveKey(passphrase, salt) {
        salt = salt || crypto.getRandomValues(new Uint8Array(8));
        return await crypto.subtle
                .importKey("raw", str2buf(passphrase), "PBKDF2", false, ["deriveKey"])
                .then(key =>
                crypto.subtle.deriveKey(
                        { name: "PBKDF2", salt, iterations: 1000, hash: "SHA-256" },
                        key,
                        { name: "AES-GCM", length: 256 },
                        false,
                        ["encrypt", "decrypt"],
                ),
        );
}

async function decrypt(passphrase, encrypted) {
        const saltIvCipherHex = encrypted;
        const [salt, iv, data] = saltIvCipherHex.split("-").map(hex2buf);
        let key = await deriveKey(passphrase, salt);
        let v = await crypto.subtle.decrypt({ name: "AES-GCM", iv }, key, data);
        return buf2str(new Uint8Array(v));
}
