import axios from "axios";

const baseURL = typeof window !== "undefined" && window.location && window.location.hostname
	? __API_URL__.replace("localhost", window.location.hostname)
	: __API_URL__;

const instance = axios.create({
	baseURL: baseURL,
	timeout: 1000 * 5
});
export default instance;