import { NavLink } from "react-router";
import { Badge } from "@/components/ui/lib/badge.tsx";
import React, { useState } from "react";
import classNames from "classnames";
import { useAuth } from "@/Auth.js";
import {
  dislikeFilm,
  dislikeFilmCouch,
  likeFilm,
  likeFilmCouch,
  removeDislike,
  removeDislikeCouch,
  removeLike,
  removeLikeCouch,
} from "@/client.js";
import { duration } from "@/duration.js";

export function interpolateColor(rating) {
  const startColor = { r: 255, g: 70, b: 70 };
  const endColor = { r: 100, g: 255, b: 100 };

  const t = Math.max(0, Math.min(1, rating / 10));

  const r = Math.round(startColor.r + t * (endColor.r - startColor.r));
  const g = Math.round(startColor.g + t * (endColor.g - startColor.g));
  const b = Math.round(startColor.b + t * (endColor.b - startColor.b));

  return `rgb(${r}, ${g}, ${b})`;
}

export function Card({ info, buttons }) {
  return (
    <div
      className={
        "bg-background w-[min(calc(100vw_-_4rem),100%)] overflow-hidden rounded-[2rem] border-[2px] border-solid border-zinc-500"
      }
    >
      <NavLink to={"/cinema/" + info.id}>
        <img
          src={info.poster_url}
          className={"h-[125px] w-full object-cover sm:h-[200px]"}
        />
      </NavLink>
      <div className={"w-full p-[1rem]"}>
        <NavLink to={"/cinema/" + info.id}>
          <h1
            className={
              "overflow-1 w-full text-center text-2xl font-extrabold sm:text-3xl"
            }
          >
            {info.name}
          </h1>
          <div
            className={"autoalt overflow-3 h-[4em] w-full text-xl sm:text-2xl"}
          >
            {info.description}
          </div>
        </NavLink>
        {info.genres && (
          <div
            className={
              "mt-[1rem] flex h-[4em] flex-row flex-wrap justify-start gap-[0.5rem] gap-y-[100px] overflow-hidden"
            }
          >
            {info.genres.map((genre, index) => (
              <NavLink key={genre} to={"/search?tags=" + genre}>
                <Badge
                  variant={"secondary"}
                  className={"button h-fit px-[0.4em] py-[0.3em] text-xl"}
                  key={index}
                >
                  {genre}
                </Badge>
              </NavLink>
            ))}
          </div>
        )}
        <div className={"flex h-[5rem] flex-row flex-nowrap justify-between"}>
          {info.duration_minutes && (
            <div className={"mt-[1rem] text-center"}>
              <Badge
                className={
                  "inline-block px-[0.4em] py-[0.3em] align-middle text-2xl"
                }
              >
                {duration(info.duration_minutes)}
              </Badge>
            </div>
          )}
          {info.rating && (
            <div className={"mt-[1rem] text-center md:text-left"}>
              <div
                className={
                  "mt-[1rem] mr-[1rem] inline text-center align-middle text-2xl md:text-left"
                }
              >
                Rating:
              </div>
              <Badge
                className={
                  "inline-block rounded-[2rem] px-[0.4em] py-[0.3em] align-middle text-2xl"
                }
                style={{ backgroundColor: interpolateColor(info.rating) }}
              >
                {Math.floor(info.rating * 10) / 10}
              </Badge>
            </div>
          )}
        </div>
        <div className={"mt-[1rem] flex flex-row flex-nowrap justify-between"}>
          {buttons}
        </div>
      </div>
    </div>
  );
}

export function CardLike({ info, group }) {
  const [liked, setLiked] = useState(info.user_like_status || 0);

  const [isAuth, , , token] = useAuth();

  function handleLike(like) {
    setLiked((x) => {
      if (x !== 0 && (x === 1) === like) return 0;
      return like ? 1 : -1;
    });
    if (!group) {
      if (liked === -1 && !like) removeDislike(info.id, token);
      if (liked === 1 && like) removeLike(info.id, token);
      if (liked !== -1 && !like) dislikeFilm(info.id, token);
      if (liked !== 1 && like) likeFilm(info.id, token);
    } else {
      if (liked === -1 && !like) removeDislikeCouch(info.id, group, token);
      if (liked === 1 && like) removeLikeCouch(info.id, group, token);
      if (liked !== -1 && !like) dislikeFilmCouch(info.id, group, token);
      if (liked !== 1 && like) likeFilmCouch(info.id, group, token);
    }
  }

  let buttons = (
    <>
      <Badge
        variant={"outline"}
        onClick={() => handleLike(false)}
        className={
          "button flex aspect-square h-[64px] w-fit rounded-[100%] px-[0.7em] text-2xl"
        }
      >
        <i
          className={classNames("bi-heartbreak-fill text-[2rem]", {
            "text-[red]": liked === -1,
            "text-primary": liked !== -1,
          })}
        ></i>
      </Badge>
      <Badge
        variant={"outline"}
        onClick={() => handleLike(true)}
        className={
          "button flex aspect-square h-[64px] w-fit rounded-[100%] px-[0.7em] text-2xl"
        }
      >
        <i
          className={classNames("bi-heart-fill text-[2rem]", {
            "text-[red]": liked === 1,
            "text-primary": liked !== 1,
          })}
        ></i>
      </Badge>
    </>
  );

  if (!isAuth) buttons = null;

  return <Card info={info} buttons={buttons} />;
}
