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
                    bsonType: "string",
                    description: "username"
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

db.createCollection("accessTokens", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["expireAt", "id", "login", "token"],
            properties: {
                expireAt: {
                    bsonType: "date",
                    description: "date of expiration"
                },
                id: {
                    bsonType: "binData",
                    description: "uuid of token pair"
                },
                login: {
                    bsonType: "string",
                    description: "username"
                },
                token: {
                    bsonType: "string",
                    description: "JWT token"
                }
            }
        }
    }
})

db.accessTokens.createIndex(
    {
        "expireAt": 1
    },
    {
        "expireAfterSeconds": 0
    }
)

db.createCollection("refreshTokens", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["expireAt", "id", "login", "token"],
            properties: {
                expireAt: {
                    bsonType: "date",
                    description: "date of expiration"
                },
                id: {
                    bsonType: "binData",
                    description: "uuid of token pair"
                },
                login: {
                    bsonType: "string",
                    description: "username"
                },
                token: {
                    bsonType: "string",
                    description: "JWT token"
                }
            }
        }
    }
})

db.refreshTokens.createIndex(
    {
        "expireAt": 1
    },
    {
        "expireAfterSeconds": 0
    }
)
