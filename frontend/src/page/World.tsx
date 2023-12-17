import { useParams } from "react-router-dom";

export function World() {
  const { worldId } = useParams();
  return <div> World {worldId} </div>;
}
