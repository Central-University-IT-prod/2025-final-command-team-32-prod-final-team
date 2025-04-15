import React from "react";
import { Badge } from "@/components/ui/lib/badge.js";

const RenderButtons = ({ filmId, onRemoveLike, onOpenDrawer }) => {
  return (
    <>
      <Badge
        onClick={() => onRemoveLike(filmId)}
        className="text-2xl h-[64px] w-fit aspect-square px-[0.7em] flex rounded-[100%] button"
      >
        <i className="bi-heartbreak-fill text-[2rem]"></i>
      </Badge>
      <Badge
        onClick={() => onOpenDrawer(filmId)}
        className="text-2xl h-[64px] w-fit aspect-square px-[0.7em] flex rounded-[100%] button"
      >
        <i className="bi-star text-[2rem]"></i>
      </Badge>
    </>
  );
};

export default RenderButtons;
