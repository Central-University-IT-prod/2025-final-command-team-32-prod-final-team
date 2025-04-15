import React, { useEffect, useState } from "react";
import { useAuth } from "@/Auth.js";
import { Navigate, NavLink, useNavigate } from "react-router";
import { getAllCouches } from "@/client.js";

export function Couches() {
  function Card({ info }) {
    return (
      <NavLink
        to={"/couches/" + info.id}
        className={
          "button bg-background h-fit overflow-hidden rounded-[2rem] border-[2px] border-solid border-zinc-500 p-[2rem]"
        }
      >
        <h1 className={"text-3xl"}>{info.name}</h1>
        <div className={"text-center text-xl text-zinc-500"}>
          {info.users.length} users
        </div>
      </NavLink>
    );
  }

  const navigate = useNavigate();

  const [couches, setCouches] = useState([]);

  useEffect(() => {
    getAllCouches(token).then((response) => setCouches(response));
  }, []);

  const [auth, , creds, token] = useAuth(true);

  if (!auth) {
    return <Navigate to="/auth" />;
  }

  return (
    <div className={"mx-[2rem]"}>
      <h1 className={"mb-[1rem] text-center text-4xl md:text-left"}>
        Your Couches
      </h1>
      <div className={"flex flex-row flex-wrap gap-[1rem]"}>
        {couches
          .filter((x) => x.author === creds.username)
          .map((x) => (
            <Card key={x.id} info={x} />
          ))}
        <div
          onClick={() => navigate("/createCouch")}
          className={
            "button flex aspect-square h-[8.25rem] w-[8.25rem] flex-col justify-center rounded-[2rem] border-[3px] border-dashed border-zinc-400"
          }
        >
          <div className="flex flex-row justify-center">
            <i className={"bi-plus h-fit w-fit text-[4rem]"} />
          </div>
        </div>
      </div>
      <h1 className={"mt-[3rem] mb-[1rem] text-center text-4xl md:text-left"}>
        Other's Couches
      </h1>
      <div className={"flex flex-row flex-wrap gap-[1rem]"}>
        {couches
          .filter((x) => x.author !== creds.username)
          .map((x) => (
            <Card key={x.id} info={x} />
          ))}
      </div>
    </div>
  );
}
