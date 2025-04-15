import React, { useState } from "react";
import { Navigate } from "react-router";
import { useAuth } from "@/Auth.tsx";
import { getCouchFeed, getFilmFeed, seen, seenCouch } from "@/client.ts";
import IntersectionObserverComponent from "@/PaginationTrigger.tsx";
import { CardLike } from "@/Pages/Explore/Card.jsx";
import { ExploreTinder } from "@/Pages/Explore/ExploreTinder.jsx";

function ExploreClassic({ data, group, nextCallback, switchCallback }) {
  return (
    <div className={"mx-auto mt-[2rem] w-[min(1500px,100%)] px-[2rem]"}>
      <h1
        className={"text-extrabold mb-[1rem] text-center text-4xl md:text-left"}
      >
        Explore
      </h1>
      <div
        onClick={() => switchCallback()}
        className={
          "button bg-primary mb-[3rem] rounded-[2rem] px-[0.6em] py-[1em] text-center text-3xl break-words text-black"
        }
      >
        Try out new Tinder-like interface!
      </div>
      <div
        className={"grid gap-[2rem]"}
        style={{
          gridTemplateColumns: "repeat(auto-fill, minmax(350px, 1fr))",
        }}
      >
        {data.map((item, index) => (
          <CardLike key={item.id} group={group} info={item} />
        ))}
        <IntersectionObserverComponent onIntersect={nextCallback} />
      </div>
    </div>
  );
}

export function Explore({ group }) {
  const auth = useAuth();

  const [tinderMode, setTinder] = React.useState(
    JSON.parse(localStorage.getItem("tinderMode")),
  );

  const [data, setData] = React.useState([]);

  const [stopped, setStopped] = useState(false);

  if (!auth[0]) {
    return <Navigate to="/auth" />;
  }

  function next() {
    if (stopped) return;
    if (!group) {
      return getFilmFeed(10, auth[3]).then((x) => {
        if (x.length === 0) {
          setStopped(true);
        }
        setData((y) => [...y, ...x]);
        seen(
          x.map((x) => x.id),
          auth[3],
        );
      });
    } else {
      return getCouchFeed(group, 10).then((x) => {
        if (x.length === 0) {
          setStopped(true);
        }
        setData((y) => [...y, ...x]);
        seenCouch(
          x.map((x) => x.id),
          group,
        );
      });
    }
  }

  function setTinderMode(status) {
    setTinder(status);
    localStorage.setItem("tinderMode", JSON.stringify(status));
  }

  return (
    <>
      {tinderMode ? (
        <ExploreTinder
          data={data}
          group={group}
          nextCallback={next}
          switchCallback={() => setTinderMode(false)}
        />
      ) : (
        <ExploreClassic
          data={data}
          group={group}
          nextCallback={next}
          switchCallback={() => setTinderMode(true)}
        />
      )}
    </>
  );
}
