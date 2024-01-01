import { Card } from "@/components/Card.tsx";
import { useParams } from "react-router-dom";
import { World } from "@/util/types.ts";
import { WorldsContext } from "@/util/WorldsContext.tsx";
import { useContext, useEffect, useState } from "react";

export function WorldPage() {
  const { worldId } = useParams();
  const [world, setWorld] = useState<World | null>(null);
  const [worlds] = useContext(WorldsContext)!;
  useEffect(() => {
    worlds.forEach((w) => {
      if (w.id === +(worldId ?? 0)) {
        setWorld(w);
      }
    });
  }, []);
  return (
    <Card h2={world?.name}>
      <></>
    </Card>
  );
}
