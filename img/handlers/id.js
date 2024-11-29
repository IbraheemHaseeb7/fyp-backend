/*
 * The purpose of this handler is to keep generating unique ids
 * everytime for a new image
 *
 */

const {v4} = require("uuid");

function generateId() {
	return v4(); 
}

module.exports = generateId
