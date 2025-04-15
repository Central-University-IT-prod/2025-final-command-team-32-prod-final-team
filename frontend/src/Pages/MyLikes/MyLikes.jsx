import React, { useEffect, useState } from "react";
import { useDevice } from "@/useDevice.js";
import { Card } from "@/Pages/Explore/Card.jsx";
import { getPlans, rateFilm, removeLike } from "@/client.js";
import { useAuth } from "@/Auth";
import { toast } from "sonner";
import RenderButtons from "@/Pages/MyLikes/RenderButtons.jsx";
import RatingDrawer from "./Rating.jsx";
import { Navigate } from "react-router";

const MyLikes = () => {
  const [auth, , , token] = useAuth(true);
  const [open, setOpen] = useState(false);
  const [selectedFilmId, setSelectedFilmId] = useState(null);
  const mobile = useDevice();
  const star = mobile ? 5 : 10;
  const [likedFilms, setLikedFilms] = useState([]);

  if (!auth) {
    return <Navigate to="/auth" />;
  }

  useEffect(() => {
    async function loadLikedFilms() {
      try {
        const films = await getPlans(token, 10);
        setLikedFilms(films);
      } catch (error) {
        toast.error("Error loading liked films: " + error.message);
      }
    }

    if (auth) {
      loadLikedFilms();
    }
  }, [auth]);

  const handleRemoveLike = async (filmId) => {
    try {
      setOpen(false);
      setLikedFilms((prev) => prev.filter((film) => film.id !== filmId));
      await removeLike(filmId, token);
    } catch (error) {
      toast.error("Error removing like: " + error.message);
    }
  };

  const handleOpenDrawer = (filmId) => {
    setSelectedFilmId(filmId);
    setOpen(true);
  };

  const handleSaveRating = async (filmId, rating) => {
    if (!filmId) return;
    try {
      handleRemoveLike(filmId);
      await rateFilm(filmId, rating, token);
    } catch (error) {
      toast.error("Error rating film: " + error.message);
    } finally {
      setSelectedFilmId(null);
      setOpen(false);
    }
  };

  return (
    <div style={{ display: "flex", flexDirection: "column" }}>
      <h1
        style={{ textAlign: "center" }}
        className={"text-extrabold mb-[1rem] text-center text-4xl md:text-left"}
      >
        My Likes
      </h1>
      <div
        className="grid gap-[2rem] px-[2rem] pt-10"
        style={{
          gridTemplateColumns: "repeat(auto-fill, minmax(350px, 1fr))",
        }}
      >
        {likedFilms.map((film) => (
          <div key={film.id}>
            <Card
              info={film}
              buttons={
                <RenderButtons
                  filmId={film.id}
                  onRemoveLike={handleRemoveLike}
                  onOpenDrawer={handleOpenDrawer}
                />
              }
            />
          </div>
        ))}
        <RatingDrawer
          open={open}
          onOpenChange={setOpen}
          star={star}
          onSaveRating={(rating) =>
            handleSaveRating(selectedFilmId, star === 5 ? rating * 2 : rating)
          }
        />
      </div>
    </div>
  );
};

export default MyLikes;
