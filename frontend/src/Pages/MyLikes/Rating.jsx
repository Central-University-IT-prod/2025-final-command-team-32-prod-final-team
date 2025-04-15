import React, { useState } from "react";
import { Drawer, DrawerContent, DrawerHeader, DrawerTitle } from "@/components/ui/lib/drawer";

const RatingDrawer = ({ open, onOpenChange, onSaveRating, star = 10 }) => {
  const [rating, setRating] = useState(null);

  const handleSave = () => {
    if (rating === null) {
      onOpenChange(false);
      return;
    }
    onSaveRating(rating);
    setRating(null);
  };

  return (
    <Drawer open={open} onOpenChange={onOpenChange}>
      <DrawerContent>
        <DrawerHeader>
          <DrawerTitle style={{ fontSize: "32px", textAlign: "center" }} className="mt-6">
            Please rate this movie.
          </DrawerTitle>
        </DrawerHeader>
        <div className="p-4 m-7" style={{ display: "flex", flexDirection: "column", alignItems: "center" }}>
          <p>Please rate this movie on a scale of 1 to {star}</p>
          <div className="flex gap-[1rem] mt-3 mb-10">
            {[...Array(star)].map((_, i) => {
              const num = i + 1;
              return (
                <button key={num} className="p-0">
                  <i
                    onClick={() => setRating(num)}
                    style={{ fontSize: "2rem" }}
                    className={`bi ${rating >= num ? "bi-star-fill" : "bi-star"} button text-primary`}
                  ></i>
                </button>
              );
            })}
          </div>
          <div className="mt-4">
            <button style={{ fontSize: "32px" }} className="button" onClick={handleSave}>
              Save
            </button>
          </div>
        </div>
      </DrawerContent>
    </Drawer>
  );
};

export default RatingDrawer;
