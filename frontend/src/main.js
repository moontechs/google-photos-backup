import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import db from './database/index'

const app = createApp(App)

app.use(router)

app.provide('db', db)

app.mount('#app')

if (await db.status() === false) {
    router.push({name: 'error500', query: {message: "Error connecting to a database. Please, check that it's properly setup and environment variables are provided. \n VITE_DB_URL, VITE_DB_NAMESPACE, VITE_DB_DATABASE"}})
}
    // await db.connect(import.meta.env.VITE_DB_URL)
    // await db.use({
    //     namespace: import.meta.env.VITE_DB_NAMESPACE,
    //     database: import.meta.env.VITE_DB_DATABASE
    // });
    // await db.signin({
    //     username: "michael@kozii.de",
    //     password: "change123",
    //     scope: "users",
    //     namespace: import.meta.env.VITE_DB_NAMESPACE,
    //     database: import.meta.env.VITE_DB_DATABASE,
    // })
