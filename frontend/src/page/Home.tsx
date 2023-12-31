import { ColumnDef } from "@tanstack/react-table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useEffect, useState } from "react";
import { Card } from "@/components/Card.tsx";
import { get } from "@/util/request.ts";
import { Server } from "@/util/types.ts";
import { DataTable } from "@/components/DataTable.tsx";

export function Home() {
  const [worlds, setWorlds] = useState([]);
  return (
    <Card>
      <Tabs defaultValue="account" className="w-[400px]">
        <TabsList className={"flex justify-items-center"}>
          <TabsTrigger value="servers">Servers</TabsTrigger>
          <TabsTrigger value="worlds">Worlds</TabsTrigger>
        </TabsList>
        <TabsContent value="servers">
          <ServerList />
        </TabsContent>
        <TabsContent value="worlds">Change your password here.</TabsContent>
      </Tabs>
    </Card>
  );
}

function ServerList() {
  const columns: ColumnDef<Server>[] = [
    {
      accessorKey: "id",
      header: "ID",
    },
    {
      accessorKey: "version",
      header: "Version",
    },
  ];
  const [servers, setServers] = useState([]);
  useEffect(() => {
    get("/servers/list").then((res) => {
      setServers(res);
    });
  }, []);
  return <DataTable columns={columns} data={servers}></DataTable>;
}
