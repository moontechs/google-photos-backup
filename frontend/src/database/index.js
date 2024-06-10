import axios from "axios";

const db= {}

axios.defaults.headers.common['Accept'] = 'application/json';
axios.defaults.headers.common['Content-Type'] = 'application/json';

db.status = async () => {
    try {
        const response = await axios.get(import.meta.env.VITE_DB_URL + "/status")

        return response.status === 200
    } catch (error) {
        console.error(error.message);
    }

    return false
}

db.signin = async (email, password) => {
    try {
        const response = await axios.post(import.meta.env.VITE_DB_URL + "/signin", {
            email: email,
            password: password,
            sc: "users",
            ns: import.meta.env.VITE_DB_NAMESPACE,
            db: import.meta.env.VITE_DB_DATABASE,
        })

        if (response.status === 200) {
            localStorage.setItem('token', response.data.token)
        }

        return true
    } catch (error) {
        localStorage.removeItem('token')
        console.error(error);
    }

    return false
}

db.signout = () => {
    localStorage.removeItem('token')
}

db.authenticated = () => {
    return localStorage.getItem('token') !== null
}

export default db