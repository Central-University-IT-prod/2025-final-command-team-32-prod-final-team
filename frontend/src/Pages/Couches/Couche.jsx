import React, { useEffect, useState } from "react";
import { useAuth } from "@/Auth.js";
import { Navigate, NavLink, useParams } from "react-router";
import { getCouch } from "@/client.js";
import { Button } from "@/components/ui/lib/button.js";

export function Couche() {
  const [info, setInfo] = useState({
    name: "Name loading...",
    users: ["user #1", "user #2"],
  });

  const params = useParams();
  const id = params.id;

  const [auth, , , token] = useAuth(true);

  useEffect(() => {
    if (auth) {
      getCouch(id, token).then((x) => setInfo(x));
    }
  }, []);

  if (!auth) {
    return <Navigate to="/auth" />;
  }

  return (
    <>
      <h1 className={"text-center text-5xl"}>{info.name}</h1>
      <div className={"text-center text-3xl text-zinc-400"}>
        {info.users.join(", ")}
      </div>
      <div className={"flex-raw mt-[2rem] flex flex-nowrap justify-center"}>
        <NavLink to={"/couches/" + id + "/feed"}>
          <Button
            style={{
              height: "fit-content",
              margin: "auto",
              textAlign: "center",
            }}
            className={"text-2xl md:text-4xl"}
          >
            <i className={"bi-compass"} /> Shared Feed
          </Button>
        </NavLink>
      </div>
      <div className={"flex-raw mt-[2rem] flex flex-nowrap justify-center"}>
        <NavLink to={"/couches/" + id + "/likes"}>
          <Button
            style={{
              height: "fit-content",
              margin: "auto",
              textAlign: "center",
            }}
            className={"text-2xl md:text-4xl"}
          >
            <i className={"bi-heart"} /> Likes
          </Button>
        </NavLink>
      </div>
    </>
  );
}
