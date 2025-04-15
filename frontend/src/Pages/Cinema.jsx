import React, { useEffect, useState } from "react";
import { Badge } from "@/components/ui/lib/badge.js";
import { duration } from "@/duration.js";
import { NavLink, useParams } from "react-router";
import { getFilm } from "@/client.js";
import { useCreateAccountDialog } from "@/CreateAccountDialog.jsx";
import { interpolateColor } from "@/Pages/Explore/Card.jsx";

export function Cinema() {
  const params = useParams();
  const id = params.id;

  const [cinema, setCinema] = useState(null);

  useEffect(() => {
    getFilm(id).then((data) => setCinema(data));
  }, []);

  const authAlert = useCreateAccountDialog();

  function changeRate() {
    authAlert();
  }

  if (cinema === null)
    return <h1 className={"text-center text-5xl"}>Loading...</h1>;

  return (
    <div className={"m-auto w-[min(1200px,100%)]"}>
      <img
        src={cinema.poster_url || "https://placehold.co/1920x1080"}
        className={"m-auto h-[200px] w-full object-cover md:h-[300px]"}
      ></img>
      <div className={"p-[1rem]"}>
        <h1
          className={
            "mt-[1.5rem] text-center text-3xl font-extrabold md:text-left md:text-4xl"
          }
        >
          {cinema.name} {cinema.year && `(${cinema.year})`}
        </h1>
        {cinema.name_original && (
          <h1
            className={
              "text-center text-2xl text-zinc-500 md:text-left md:text-3xl"
            }
          >
            {cinema.name_original}
          </h1>
        )}
        {cinema.rating && (
          <div className={"mt-[1rem] text-center md:text-left"}>
            <div
              className={
                "mt-[1rem] mr-[1rem] inline text-center align-middle text-2xl md:text-left"
              }
            >
              Rating:
            </div>
            <Badge
              style={{ backgroundColor: interpolateColor(cinema.rating) }}
              className={
                "inline-block rounded-[2rem] px-[0.4em] py-[0.3em] align-middle text-2xl"
              }
            >
              {Math.floor(cinema.rating * 10) / 10}
            </Badge>
          </div>
        )}
        {/*<div className={"mt-[1rem] text-center md:text-left"}>
                <div className={"mt-[1rem] align-middle text-2xl inline text-center md:text-left mr-[1rem]"}>
                    Your rate:
                </div>
                <div className={"inline-block align-middle"}>
                    <Badge onClick={()=>changeRate()}
                        className={"text-2xl w-fit aspect-square px-[0.7em] flex rounded-[100%] button"}>{cinema.user_rate || "+"}</Badge>
                </div>
            </div>*/}
        {cinema.genres && (
          <div
            className={
              "mt-[1rem] flex flex-row flex-wrap justify-center gap-[0.5rem] md:justify-start"
            }
          >
            {cinema.genres.map((genre, index) => (
              <NavLink to={"/search?tags=" + genre}>
                <Badge
                  variant={"secondary"}
                  className={"button px-[0.4em] py-[0.3em] text-2xl"}
                  id={index}
                >
                  {genre}
                </Badge>
              </NavLink>
            ))}
          </div>
        )}
        {cinema.description && (
          <div className={"mt-[1rem] text-center text-2xl md:text-left"}>
            {cinema.description}
          </div>
        )}
        {cinema.age_rating && (
          <div className={"mt-[2rem] text-center md:text-left"}>
            <div
              className={
                "mt-[1rem] mr-[1rem] inline text-center align-middle text-2xl md:text-left"
              }
            >
              Age:
            </div>
            <Badge
              className={
                "inline-block px-[0.4em] py-[0.3em] align-middle text-2xl"
              }
            >
              {cinema.age_rating}+
            </Badge>
          </div>
        )}
        {cinema.duration_minutes && (
          <div className={"mt-[1rem] text-center md:text-left"}>
            <div
              className={
                "mt-[1rem] mr-[1rem] inline text-center align-middle text-2xl md:text-left"
              }
            >
              Duration:
            </div>
            <Badge
              className={
                "inline-block px-[0.4em] py-[0.3em] align-middle text-2xl"
              }
            >
              {duration(cinema.duration_minutes)}
            </Badge>
          </div>
        )}
        {cinema.actors && (
          <>
            <h1 className={"mt-[2rem] mb-[0.5rem] text-2xl"}>Actors: </h1>
            <div
              className={
                "flex flex-row flex-wrap justify-center gap-[0.5rem] md:justify-start"
              }
            >
              {cinema.actors.map((genre, index) => (
                <Badge
                  variant={"secondary"}
                  className={"px-[0.4em] py-[0.3em] text-2xl"}
                  id={index}
                >
                  {genre}
                </Badge>
              ))}
            </div>
          </>
        )}

        <NavLink
          to={`https://www.kinopoisk.ru/index.php?kp_query=${cinema.name}`}
        >
          <div
            className={
              "button text-primary bg-secondary mt-[2rem] rounded-[1rem] p-[1rem] text-2xl"
            }
          >
            <i className={"bi-tv-fill mr-[1rem]"} />
            <div className={"inline align-middle"}>
              Смотреть на <span className={"text-[orange]"}>Кинопоиск</span>
            </div>
          </div>
        </NavLink>

        {/*<div className={"flex mt-[4rem] mb-[4rem] flex-row flex-nowrap justify-between mt-[1rem] px-[3rem]"}>
                <Badge
                       className={"text-4xl h-[5rem] bg-red-500 w-fit aspect-square px-[0.7em] flex rounded-[100%] button"}>
                    <i className={"bi-heartbreak-fill text-[black] text-[2rem]"}></i>
                </Badge>
                <Badge
                       className={"text-4xl h-[5rem] bg-green-500 w-fit aspect-square px-[0.7em] flex rounded-[100%] button"}>
                    <i className={"bi-heart-fill text-[black] text-[2rem]"}></i>
                </Badge>
            </div>*/}
      </div>
    </div>
  );
}
