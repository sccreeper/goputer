import Dexie from "dexie"

const databaseName = "goputerEditor"
const databaseVersion = 1
const fileTableName = "files"

const db = new Dexie(databaseName);

db.version(databaseVersion).stores({
    files: "name, data, type"
})

export {db, fileTableName, databaseVersion}

