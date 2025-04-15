import faker
import requests
from faker.generator import random
from faker.proxy import Faker

url = "http://prod-team-32-n26k57br.REDACTED:8080"
#create admin user
r = requests.post(url+"/api/v1/users/sign-up", json={"username":"admin", "password":"admin"})
token = r.json()["token"]

genres = [
    "Action",
    "Adventure",
    "Animation",
    "Anime",
    "Biographical",
    "Comedy",
    "Crime",
    "Drama",
    "Fantasy",
    "Historical",
    "Horror",
    "Mystery",
    "Romance",
    "Science-Fiction",
    "Thriller",
    "Western",
    "War",
    "Musical",
    "Documentary",
    "Noir",
    "Heist",
    "Martial-Arts",
    "Melodrama",
    "Neo-Noir",
    "Psychological-Thriller",
    "Road-Movie",
    "Satire",
    "Slapstick",
    "Social-Commentary",
    "Sports",
    "Spy",
    "Steampunk",
    "Superhero",
    "Surrealist",
    "Teen-Comedy",
    "Tragedy",
    "Vampire",
    "War-Drama",
    "Zombie-Apocalypse",
    "Cyberpunk",
    "Dark-Comedy",
    "Disaster",
    "Epic",
    "Post-Apocalyptic",
    "Detective",
    "Sitcom",
]


fake = Faker()

for i in range(100):
    name = fake.text(max_nb_chars=20)
    release_year = fake.random_int(min=1900, max=2025)
    age_rating = fake.random_int(min=0, max=18)
    duration_minutes = fake.random_int(min=1, max=360)
    poster_url = "https://random-image-pepebigotes.vercel.app/api/random-image"
    description = fake.text(max_nb_chars=100)
    req_genres = random.sample(genres, random.randint(1, 5))
    Actors = [
        fake.name(),
        fake.name(),
        fake.name(),
    ]
    print(len(req_genres))
    rating= float(fake.random_int(min=10, max=100)/10.0)
    req = {
        "name": name,
        "year": release_year,
        "age_rating": age_rating,
        "description": description,
        "genres": req_genres,
        "actors": Actors,
        "rating": rating,
        "user_rating": fake.random_int(min=0, max=5),
        "duration_minutes": duration_minutes,
        "original_title": fake.text(max_nb_chars=20),
        "poster_url" : poster_url
    }
    res = requests.post(url+"/api/v1/films", json=req, headers={"Authorization": "Bearer " + token})
    print(res.text)
# create cinema

# Name            string   `json:"name" validate:"required"`
# ReleaseYear     *int32   `json:"year,omitempty"`
# AgeRating       *int32   `json:"age_rating,omitempty"`
# DurationMinutes *int32   `json:"duration_minutes,omitempty"`
# PosterURL       string   `json:"poster_url,omitempty"`
# Description     string   `json:"description,omitempty"`
# Genres          []string `json:"genres,omitempty"`
# Actors          []string `json:"actors,omitempty"`
# Rating          float32  `json:"rating,omitempty"`
# UserRating      int      `json:"user_rating,omitempty"`
# Private         bool     `json:"private"`
# OriginalTitle   *string  `json:"original_title,omitempty"`

