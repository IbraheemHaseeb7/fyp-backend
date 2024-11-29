/*
 * stores images on locally against a unique id
 */
const fs = require("fs");
const path = require("path");
const generateId = require("./id")

function store(image) {
	const dir = path.join(__dirname, "..", "uploads", "files");

	if (!fs.existsSync(dir)) {
		fs.mkdirSync(dir);
	}

	const imagePath = generateId();
	const filePath = `${path.join(dir, imagePath)}.${image.mimetype.split("/")[1]}`;

	fs.writeFile(filePath, image.buffer, err => {
		if (err) {
			throw new Error("File could not be saved");
		}
	})

	return `${imagePath}.${image.mimetype.split("/")[1]}`;
}

module.exports = {
	store
}
