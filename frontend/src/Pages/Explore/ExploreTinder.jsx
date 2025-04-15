import React, { useEffect, useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { Card } from "@/Pages/Explore/Card.jsx";
import { Badge } from "@/components/ui/lib/badge.tsx";
import { useAuth } from "@/Auth.js";
import {
  dislikeFilm,
  dislikeFilmCouch,
  likeFilm,
  likeFilmCouch,
} from "@/client.js";

export function ExploreTinder({ data, group, nextCallback, switchCallback }) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [showCard, setShowCard] = useState(true);
  const [direction, setDirection] = useState(0);

  const cardVariants = {
    initial: { x: 0, opacity: 1, rotate: 0 },
    animate: { x: 0, opacity: 1, rotate: 0 },
    exit: {
      x: direction * window.innerWidth,
      opacity: 0,
      rotate: direction * 30,
      transition: { duration: 0.5, ease: "easeInOut" },
    },
  };

  const [isAuth, , , token] = useAuth();

  const [canLike, setCanLike] = useState(true);

  useEffect(() => {
    setCanLike(true);
  }, [currentIndex]);

  function handleLike(like) {
    if (!canLike) return;
    setCanLike(false);
    if (!group) {
      if (!like) dislikeFilm(data[currentIndex].id, token);
      else likeFilm(data[currentIndex].id, token);
    } else {
      if (!like) dislikeFilmCouch(data[currentIndex].id, group, token);
      else likeFilmCouch(data[currentIndex].id, group, token);
    }
  }

  useEffect(() => {
    if (currentIndex >= data.length - 3) {
      nextCallback();
    }
  }, [currentIndex, data.length, nextCallback]);

  useEffect(() => {
    console.log("direction", direction);
    if (direction !== 0) {
      setShowCard(false);
    }
  }, [direction]);

  function swipe(like) {
    handleLike(like);
    setDirection(like ? 1 : -1);
  }

  return (
    <div className="h-fit w-full overflow-x-hidden">
      <div className="mx-auto w-[min(calc(100%-4rem),500px)]">
        <div className="relative flex h-[65vh] w-full items-center justify-center px-[2rem]">
          {data[currentIndex + 1] && (
            <div className="absolute z-0 w-full">
              <Card info={data[currentIndex + 1]} />
            </div>
          )}

          <AnimatePresence
            onExitComplete={() => {
              setCurrentIndex((prev) => prev + 1);
              setDirection(0);
              setShowCard(true);
            }}
          >
            {showCard && data[currentIndex] && (
              <motion.div
                key={currentIndex}
                className="absolute z-10 w-full"
                variants={cardVariants}
                initial="initial"
                animate="animate"
                exit="exit"
              >
                <Card info={data[currentIndex]} />
              </motion.div>
            )}
          </AnimatePresence>
        </div>

        <div className="z-[8] mt-[1rem] flex flex-row flex-nowrap justify-between px-[3rem]">
          <Badge
            onClick={() => swipe(false)}
            className="button flex aspect-square h-[5rem] w-fit rounded-full bg-red-500 px-[0.7em] text-4xl"
          >
            <i className="bi-heartbreak-fill text-3xl text-black"></i>
          </Badge>
          <Badge
            onClick={() => swipe(true)}
            className="button flex aspect-square h-[5rem] w-fit rounded-full bg-green-500 px-[0.7em] text-4xl"
          >
            <i className="bi-heart-fill text-3xl text-black"></i>
          </Badge>
        </div>
        <div
          className="button mt-[1rem] mb-[1rem] text-center text-2xl underline"
          onClick={switchCallback}
        >
          Go back to classic mode 🙁
        </div>
      </div>
    </div>
  );
}
