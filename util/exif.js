const fs = require("fs");
((packname = "Hello👋", author = "Fdvky") => {
  const json = {
    "sticker-pack-id": "fdb.my.id 1456",
    "sticker-pack-name": packname,
    "sticker-pack-publisher": author,
    "android-app-store-link":
      "https://play.google.com/store/apps/details?id=com.rayark.cytus2",
    "ios-app-store-link": "https://apps.apple.com/app/id625334537",
    emojis: ["👋"],
  };
  const littleEndian = Buffer.from([
    0x49, 0x49, 0x2a, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01, 0x00, 0x41, 0x57,
    0x07, 0x00,
  ]);
  const bytes = [0x00, 0x00, 0x16, 0x00, 0x00, 0x00];

  let len = new TextEncoder().encode(JSON.stringify(json)).length;
  let last;

  if (len > 256) {
    len = len - 256;
    bytes.unshift(0x01);
  } else {
    bytes.unshift(0x00);
  }

  if (len < 16) {
    last = len.toString(16);
    last = "0" + len;
  } else {
    last = len.toString(16);
  }

  const buf2 = Buffer.from(last, "hex");
  const buf3 = Buffer.from(bytes);
  const buf4 = Buffer.from(JSON.stringify(json));

  const buffer = Buffer.concat([littleEndian, buf2, buf3, buf4]);
  fs.writeFileSync("./raw.exif", buffer);
  process.exit(0);
})();
