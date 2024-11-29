// this file is just to run the application locally as  a temp server
const dotenv = require("dotenv");
const express = require("express");
const { store } = require("./handlers/images");
const multer = require("multer");
const app = express();
const path = require("path")

// loading env variables
dotenv.config()

// setting up multer middleware
const upload = multer({store:"uploads/"})

app.post("/upload", upload.single("image"), (req, res) => {
	const uri = store(req.file);
	return res.json({message: "successfully stored image", uri})
})

app.get("/get/:id", (req, res) => {
	return res.json({message: "image"});
})

app.use(express.static(path.join(__dirname, 'uploads')));

app.listen(3000, _=>console.log("Listening on port: 3000"))
