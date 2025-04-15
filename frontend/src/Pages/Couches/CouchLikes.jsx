import React, { useEffect, useState } from "react";
import { useDevice } from "@/useDevice.js";
import { Card } from "@/Pages/Explore/Card.jsx";
import { getCouchPlans, removeLikeCouch } from "@/client.js";
import { useAuth } from "@/Auth";
import { toast } from "sonner";
import RenderButtons from "@/Pages/MyLikes/RenderButtons.jsx";
import { Navigate, useParams } from "react-router-dom";
import RatingDrawer from "@/Pages/MyLikes/Rating.jsx";

const CouchLikes = () => {
  const { id: couchId } = useParams();
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
    async function loadCouchFilms() {
      try {
        const films = await getCouchPlans(couchId, token, 30);
        setLikedFilms(films);
      } catch (error) {
        toast.error("Error loading couch films: " + error.message);
      }
    }

    if (auth) {
      loadCouchFilms();
    }
  }, [auth, couchId, token]);

  const handleRemoveLike = async (filmId) => {
    try {
      setLikedFilms((prev) => prev.filter((film) => film.id !== filmId));
      await removeLikeCouch(filmId, couchId, token);
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
    } catch (error) {
      toast.error("Error rating film: " + error.message);
    } finally {
      setSelectedFilmId(null);
      setOpen(false);
    }
  };

  return (
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
  );
};

export default CouchLikes;
