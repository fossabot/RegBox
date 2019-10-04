db = connect("localhost:27017/admin")
db.auth('root', 'P@ssw0rd')

db = db.getSiblingDB("regbox")

db.createUser({
    user: "regbox",
    pwd: "P@ssw0rd",
    roles: ["readWrite"],

})

db.createCollection("creds", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["login", "passwd", "salt"],
            properties: {
                login: {
                    bsonType: "binData",
                    description: "hash of login, required as bin"
                },
                passwd: {
                    bsonType: "binData",
                    description: "argon2 hash of password, required as bin"
                },
                salt: {
                    bsonType: "binData",
                    description: "salt used in password, required as bin"
                }
            }
        }
    }
})

db.creds.createIndex(
    {
        "login": 1
    },
    {
        "name": "uniqueLoginIndex",
        "unique": true
    }
)
