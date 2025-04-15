import * as React from "react";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { getFilm, getGenres, searchFilm } from "@/client.js";
import { useAuth } from "@/Auth.js";
import { useNavigate } from "react-router-dom";
import { Input } from "@/components/ui/lib/input";
import { Button } from "@/components/ui/button";

const FilmSearchList = ({
  selectedFilm,
  setSelectedFilm,
  filmQuery,
  setFilmQuery,
}) => {
  const [foundFilms, setFoundFilms] = useState([]);
  const wrapperRef = useRef(null);

  // Выполняем поиск фильмов с дебаунсом
  useEffect(() => {
    if (filmQuery.length > 2) {
      const timeoutId = setTimeout(() => {
        searchFilm(filmQuery)
          .then(setFoundFilms)
          .catch((error) => {
            console.error(error);
            toast.error("Error searching films");
          });
      }, 500);
      return () => clearTimeout(timeoutId);
    } else {
      setFoundFilms([]);
    }
  }, [filmQuery]);

  // Закрытие выпадающего списка при клике вне его области
  useEffect(() => {
    function handleClickOutside(event) {
      if (wrapperRef.current && !wrapperRef.current.contains(event.target)) {
        setFoundFilms([]);
      }
    }

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div className="relative" ref={wrapperRef}>
      <Input
        type="text"
        value={filmQuery}
        onChange={(e) => setFilmQuery(e.target.value)}
        placeholder="Enter film name..."
        className="w-full"
      />
      {foundFilms.length > 0 && (
        <ul className="absolute right-0 left-0 z-10 mt-1 rounded-md border bg-black shadow-lg">
          {foundFilms.map((film) => (
            <li
              key={film.id}
              onClick={() => {
                setSelectedFilm(film);
                setFilmQuery(film.name);
                setFoundFilms([]);
              }}
              className="cursor-pointer p-2 hover:bg-[#262626]"
            >
              {film.name}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

const AdminEditFilm = () => {
  const navigate = useNavigate();
  const [, , , token] = useAuth();
  const [selectedFilm, setSelectedFilm] = useState(null);
  const [filmQuery, setFilmQuery] = useState("");
  const [filmName, setFilmName] = useState("");
  const [description, setDescription] = useState("");
  const [genres, setGenres] = useState([]); // Массив строк
  const [selectedFile, setSelectedFile] = useState(null);
  const [posterUrl, setPosterUrl] = useState("");
  const [loading, setLoading] = useState(false);

  // Загрузка списка жанров
  useEffect(() => {
    getGenres()
      .then((genres) => {
        // Здесь genres – массив строк
        // При необходимости можно создать дополнительные опции для селектора
      })
      .catch((error) => {
        console.error(error);
        toast.error("Error loading genres");
      });
  }, []);

  // При выборе фильма загружаем его данные
  useEffect(() => {
    if (selectedFilm) {
      getFilm(selectedFilm.id)
        .then((film) => {
          setFilmName(film.name);
          setDescription(film.description || "");
          setPosterUrl(film.poster_url || "");
          setGenres(film.genres || []);
        })
        .catch((error) => {
          console.error(error);
          toast.error("Error loading film data");
        });
    }
  }, [selectedFilm]);

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!filmName) {
      toast.error("Please provide a film name");
      return;
    }

    setLoading(true);

    const formData = new FormData();
    formData.append("film_name", filmName);
    formData.append("description", description);
    formData.append("genres", JSON.stringify(genres)); // Жанры передаются как JSON‑строка
    if (selectedFile) {
      formData.append("uploadfile", selectedFile);
    }

    try {
      const response = await fetch(
        `https://prod-team-32-n26k57br.REDACTED/api/v1/films/${selectedFilm.id}`,
        {
          method: "PUT",
          headers: {
            Authorization: `Bearer ${token}`,
          },
          body: formData,
        },
      );

      if (!response.ok) {
        const errorText = await response.text();
        toast.error(`Error: ${errorText}`);
        return;
      }

      toast.success("Film updated successfully!");
      navigate("/adminPanel");
    } catch (error) {
      console.error(error);
      toast.error("Error updating film");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className="mx-auto max-w-xl space-y-4 p-4"
      style={{ minWidth: "45vw" }}
    >
      <h1 className="mb-4 text-3xl">Edit Film</h1>
      <div className="space-y-2">
        <label className="block">
          <p className="mb-2 text-xl">Search Film</p>
          <FilmSearchList
            selectedFilm={selectedFilm}
            setSelectedFilm={setSelectedFilm}
            filmQuery={filmQuery}
            setFilmQuery={setFilmQuery}
          />
        </label>
        {selectedFilm && (
          <form onSubmit={handleSubmit} className="space-y-4">
            <label className="block">
              <p className="mb-2 text-xl">Film Title</p>
              <Input
                value={filmName}
                onChange={(e) => setFilmName(e.target.value)}
                placeholder="Film title"
              />
            </label>
            <label className="block">
              <p className="mb-2 text-xl">Description</p>
              <Input
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Film description"
              />
            </label>
            <label className="block">
              <p className="mb-2 text-xl">Genres (comma separated)</p>
              <Input
                value={genres.join(",")}
                onChange={(e) =>
                  setGenres(e.target.value.split(",").map((g) => g.trim()))
                }
                placeholder="Enter genres"
              />
            </label>
            <label className="block">
              <p className="mb-2 text-xl">Poster</p>
              {posterUrl && (
                <img
                  src={posterUrl}
                  alt="Poster preview"
                  className="mt-2 h-32 object-contain"
                />
              )}
            </label>
            <Button type="submit" disabled={loading} className="w-full text-lg">
              {loading ? "Updating..." : "Save Changes"}
            </Button>
          </form>
        )}
      </div>
    </div>
  );
};

export default AdminEditFilm;
