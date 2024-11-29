const compress_images = require("compress-images");
const path = require("path")
import imagemin from 'imagemin';
import imageminJpegtran from 'imagemin-jpegtran';
import imageminPngquant from 'imagemin-pngquant';

function compressImage(image_path) {
	return new Promise((resolve, reject) => {

		//let new_image_path = path.join(__dirname, "uploads", "files", image_path);
		imagemin(["uploads/files/*.{png}"], {
			destination: "build/images",
			plugins: [
				imageminJpegtran(),
				imageminPngquant({
					quality: [0.6, 0.8]
				})
			]
		})
	});
}

module.exports = {
	compressImage
}
