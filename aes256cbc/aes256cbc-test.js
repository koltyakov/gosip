const crypto = require('crypto');

const ALGORITHM = 'aes256'; // 'aes-256-cbc';
const CIPHER_KEY = "abcdefghijklmnopqrstuvwxyz012345";  // Same key used in Golang
const BLOCK_SIZE = 16; // 16;

// Encrypts plain text into cipher text
function encrypt(plainText) {
  const iv = crypto.randomBytes(BLOCK_SIZE);
  const cipher = crypto.createCipheriv(ALGORITHM, CIPHER_KEY, iv);
  let cipherText;
  try {
    cipherText = iv.toString('hex') + cipher.update(plainText, 'utf8', 'hex') + cipher.final('hex');
  } catch (e) {
    cipherText = null;
  }
  return cipherText;
}

// Decrypts cipher text into plain text
function decrypt(cipherText) {
  const contents = Buffer.from(cipherText, 'hex');
  const iv = contents.slice(0, BLOCK_SIZE);
  const textBytes = contents.slice(BLOCK_SIZE);
  const decipher = crypto.createDecipheriv(ALGORITHM, CIPHER_KEY, iv);
  const decrypted = decipher.update(textBytes, 'hex', 'utf8') + decipher.final('utf8');
  return decrypted;
}

(() => {

  const original = "My string";
  const encoded = encrypt(original);
  const decoded = decrypt(encoded);
  let decodedFromGo = '';
  try {
    decodedFromGo = decrypt("b0f620d6616995c29b7fc1f2751ec1eb2515a27959ea572a93fc6568a235e7e8");
  } catch (ex) { /**/ }

  console.log({ original, encoded, decoded, decodedFromGo });

  const b = new Buffer(encoded).toString('base64');
  console.log({ b });

})();

// node ./aes256cbc/aes256cbc-test.js