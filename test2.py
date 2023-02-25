import requests

body = {
    "username": "astrog",
    "type": "t4",
    "team": [
        {
            "id": 167,
            "player": "astrog",
            "name": "Pyre",
            "color": "Red",
            "uid": "starter-167-90210",
            "card_detail_id": 167,
            "edition": 4,
            "editions": "4",
            "type": "Summoner",
            "rarity": 2,
            "delegated_to": None,
            "level": 1,
            "mana": 3,
            "stats": {
                "mana": 3,
                "attack": 0,
                "ranged": 0,
                "magic": 0,
                "armor": 0,
                "health": 0,
                "speed": 1
            },
            "selected": True,
            "filtered": False,
            "imgURL": "https://d36mxiodymuqjm.cloudfront.net/cards_by_level/untamed/Pyre_lv1.png",
            "abilities": [
                "+1 Speed"
            ],
            "starter": True,
            "cached_imgURL": "http://splinterforge.janglehost.com/image_server/Pyre_lv1.png",
            "race": "Pyre_lv1"
        },
        {
            "cid": "G7-441-IZBL1FE6C0",
            "player": "astrog",
            "name": "General Sloan",
            "color": "White",
            "uid": "G7-441-IZBL1FE6C0",
            "card_detail_id": 441,
            "id": 441,
            "editions": "7",
            "edition": 7,
            "type": "Summoner",
            "rarity": 2,
            "delegated_to": None,
            "level": 2,
            "mana": 4,
            "stats": {
                "mana": 4,
                "attack": 0,
                "ranged": 1,
                "magic": 0,
                "armor": 0,
                "health": 0,
                "speed": 0
            },
            "selected": True,
            "filtered": False,
            "imgURL": "https://d36mxiodymuqjm.cloudfront.net/cards_by_level/chaos/General Sloan_lv2.png",
            "abilities": [
                "+1 Ranged"
            ],
            "starter": False,
            "cached_imgURL": "http://splinterforge.janglehost.com/image_server/General Sloan_lv2.png",
            "race": "General Sloan_lv2"
        },
        {
            "id": 196,
            "player": "astrog",
            "name": "Tower Griffin",
            "color": "Gray",
            "uid": "starter-196-90210",
            "card_detail_id": 196,
            "edition": 4,
            "editions": "4",
            "type": "Monster",
            "rarity": 2,
            "delegated_to": None,
            "level": 1,
            "mana": 4,
            "stats": {
                "mana": [
                    4,
                    4,
                    4,
                    4,
                    4,
                    4,
                    4,
                    4
                ],
                "attack": [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                "ranged": [
                    1,
                    1,
                    1,
                    1,
                    1,
                    1,
                    1,
                    2
                ],
                "magic": [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                "armor": [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                "health": [
                    3,
                    4,
                    4,
                    5,
                    4,
                    5,
                    6,
                    6
                ],
                "speed": [
                    3,
                    3,
                    4,
                    4,
                    4,
                    4,
                    4,
                    4
                ],
                "abilities": [
                    [
                        "Flying"
                    ],
                    [],
                    [],
                    [],
                    [
                        "Protect"
                    ],
                    [],
                    [],
                    []
                ]
            },
            "selected": True,
            "filtered": False,
            "imgURL": "https://d36mxiodymuqjm.cloudfront.net/cards_untamed/Tower Griffin.png",
            "abilities": [
                "Flying"
            ],
            "starter": True,
            "attack": 0,
            "magic": 0,
            "ranged": 1,
            "armor": 0,
            "speed": 3,
            "health": 3,
            "cached_imgURL": "http://splinterforge.janglehost.com/image_server/Tower Griffin.png",
            "race": "Tower Griffin"
        },
        {
            "cid": "G4-157-1NCFITBQJ4",
            "player": "astrog",
            "name": "Kobold Bruiser",
            "color": "Red",
            "uid": "G4-157-1NCFITBQJ4",
            "card_detail_id": 157,
            "id": 157,
            "editions": "4",
            "edition": 4,
            "type": "Monster",
            "rarity": 1,
            "delegated_to": None,
            "level": 3,
            "mana": 3,
            "stats": {
                "mana": [
                    3,
                    3,
                    3,
                    3,
                    3,
                    3,
                    3,
                    3,
                    3,
                    3
                ],
                "attack": [
                    2,
                    2,
                    2,
                    3,
                    3,
                    3,
                    3,
                    3,
                    3,
                    4
                ],
                "ranged": [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                "magic": [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                "armor": [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                "health": [
                    3,
                    4,
                    4,
                    4,
                    5,
                    5,
                    6,
                    6,
                    7,
                    7
                ],
                "speed": [
                    2,
                    2,
                    3,
                    3,
                    3,
                    3,
                    3,
                    4,
                    4,
                    4
                ],
                "abilities": [
                    [],
                    [],
                    [],
                    [],
                    [],
                    [
                        "Knock Out"
                    ],
                    [],
                    [],
                    [],
                    []
                ]
            },
            "selected": True,
            "filtered": False,
            "imgURL": "https://d36mxiodymuqjm.cloudfront.net/cards_untamed/Kobold Bruiser.png",
            "abilities": [],
            "starter": False,
            "attack": 2,
            "magic": 0,
            "ranged": 0,
            "armor": 0,
            "speed": 3,
            "health": 4,
            "cached_imgURL": "http://splinterforge.janglehost.com/image_server/Kobold Bruiser.png",
            "race": "Kobold Bruiser"
        },
        {
            "cid": "G7-395-VNYV3SQC68",
            "player": "astrog",
            "name": "Radiated Scorcher",
            "color": "Red",
            "uid": "G7-395-VNYV3SQC68",
            "card_detail_id": 395,
            "id": 395,
            "editions": "7",
            "edition": 7,
            "type": "Monster",
            "rarity": 1,
            "delegated_to": None,
            "level": 3,
            "mana": 1,
            "stats": {
                "mana": [
                    1,
                    1,
                    1,
                    1,
                    1,
                    1,
                    1,
                    1,
                    1,
                    1
                ],
                "attack": [
                    1,
                    1,
                    1,
                    1,
                    1,
                    1,
                    2,
                    2,
                    2,
                    2
                ],
                "ranged": [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                "magic": [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                "armor": [
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0
                ],
                "health": [
                    2,
                    2,
                    2,
                    2,
                    2,
                    2,
                    2,
                    2,
                    2,
                    3
                ],
                "speed": [
                    2,
                    2,
                    2,
                    2,
                    3,
                    3,
                    3,
                    3,
                    4,
                    4
                ],
                "abilities": [
                    [],
                    [],
                    [
                        "Shatter"
                    ],
                    [],
                    [],
                    [],
                    [],
                    [],
                    [],
                    []
                ]
            },
            "selected": True,
            "filtered": False,
            "imgURL": "https://d36mxiodymuqjm.cloudfront.net/cards_chaos/Radiated Scorcher.jpg",
            "abilities": [
                "Shatter"
            ],
            "starter": False,
            "attack": 1,
            "magic": 0,
            "ranged": 0,
            "armor": 0,
            "speed": 2,
            "health": 2,
            "cached_imgURL": "http://splinterforge.janglehost.com/image_server/Radiated Scorcher.png",
            "race": "Radiated Scorcher"
        },
        {
            "player": "astrog",
            "name": "astrog",
            "color": "Gray",
            "uid": "Warrior hero",
            "card_detail_id": 0,
            "id": 0,
            "editions": "0",
            "edition": 0,
            "type": "Monster",
            "rarity": 4,
            "delegated_to": None,
            "level": 1,
            "mana": 0,
            "starter": False,
            "stats": {
                "mana": [
                    0,
                    0,
                    0,
                    0
                ],
                "attack": [
                    1,
                    1,
                    1,
                    1
                ],
                "ranged": [
                    0,
                    0,
                    0,
                    0
                ],
                "magic": [
                    0,
                    0,
                    0,
                    0
                ],
                "armor": [
                    1,
                    1,
                    1,
                    1
                ],
                "health": [
                    1,
                    1,
                    1,
                    1
                ],
                "speed": [
                    1,
                    1,
                    1,
                    1
                ],
                "abilities": [
                    [],
                    [],
                    [],
                    []
                ]
            },
            "selected": False,
            "filtered": False,
            "imgURL": "/assets/menu/Warrior_tn.jpg",
            "abilities": [],
            "gold": False,
            "hidden": False,
            "actualLevel": 4,
            "attack": 0,
            "ranged": 0,
            "magic": 0,
            "armor": 0,
            "health": 0,
            "roundSpeed": 0
        }
    ],
    "boosts": [],
    "deckPower": 0,
    "memo": "64226ef93ad9d318c8a09804b5fd07b6cfeeff10f373efa81d58db33348f1faa6fa9616a7ed16eab2503512aa8d14511c3b58f96c7c729eead00310e96e235768c14d0daa1f9512edb8bb00fa8294415d176c0b94b8a07a9018b8415c11e7cd8"
}


response = requests.post("https://splinterforge.io/boss/fight_boss", json=body)
print(response.text)
