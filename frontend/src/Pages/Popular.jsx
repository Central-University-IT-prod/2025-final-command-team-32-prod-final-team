import { CardLike } from "@/Pages/Explore/Card.jsx";
import IntersectionObserverComponent from "@/PaginationTrigger.js";
import React, { useState } from "react";
import { getPopularFilms } from "@/client.js";

export function Popular() {
  const [data, setData] = React.useState([]);

  const [stopped, setStopped] = useState(false);

  const [page, setPage] = useState(0);

  const pageSize = 10;

  function next() {
    if (stopped) return;
    return getPopularFilms(pageSize, page * pageSize).then((x) => {
      if (x.length === 0) {
        setStopped(true);
      }
      setData((y) => [...y, ...x]);
      setPage((x) => x + 1);
    });
  }

  return (
    <div className={"mx-auto mt-[2rem] w-[min(1500px,100%)] px-[2rem]"}>
      <h1
        className={"text-extrabold mb-[1rem] text-center text-4xl md:text-left"}
      >
        Popular
      </h1>
      <div
        className={"grid gap-[2rem]"}
        style={{
          gridTemplateColumns: "repeat(auto-fill, minmax(350px, 1fr))",
        }}
      >
        {data.map((item, index) => (
          <CardLike key={item.id} info={item} />
        ))}
        <IntersectionObserverComponent onIntersect={() => next()} />
      </div>
    </div>
  );
}
