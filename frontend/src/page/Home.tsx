import { MoreHorizontal } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ColumnDef } from "@tanstack/react-table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import React, { SetStateAction, useEffect, useState } from "react";
import { Card } from "@/components/Card.tsx";
import { del, get } from "@/util/request.ts";
import { Server, World } from "@/util/types.ts";
import { DataTable } from "@/components/DataTable.tsx";
import { Link } from "react-router-dom";

export function Home() {
  const [alertText, setAlertText] = useState("");
  return (
    <Card alertText={alertText}>
      <Tabs defaultValue="worlds" className={"w-11/12"}>
        <TabsList className={"flex justify-around"}>
          <TabsTrigger value="worlds" className={"w-1/2"}>
            Worlds
          </TabsTrigger>
          <TabsTrigger value="servers" className={"w-1/2"}>
            Servers
          </TabsTrigger>
        </TabsList>
        <TabsContent value="servers">
          <ServerList setAlertText={setAlertText} />
        </TabsContent>
        <TabsContent value="worlds">
          <WorldList setAlertText={setAlertText} />
        </TabsContent>
      </Tabs>
    </Card>
  );
}

function WorldList({
  setAlertText,
}: {
  setAlertText: React.Dispatch<SetStateAction<string>>;
}) {
  const columns: ColumnDef<World>[] = [
    {
      accessorKey: "id",
      header: "ID",
    },
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      id: "actions",
      cell: ({ row }) => {
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <Link to={`/worlds/${row.original.id}`}>
                <DropdownMenuItem>View</DropdownMenuItem>
              </Link>
              <DropdownMenuItem
                className={"red"}
                onClick={() => {
                  deleteWorld(row.original.id);
                }}
              >
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];
  const [worlds, setWorlds] = useState<World[]>([]);
  const [refresh, setRefresh] = useState(1);

  function deleteWorld(id: number) {
    del(`/worlds/delete`, { id: id }).then((res) => {
      setAlertText(res);
      setRefresh((prev) => prev + 1);
    });
  }

  useEffect(() => {
    setAlertText("");
    get("/worlds/list").then((res) => {
      setWorlds(res);
    });
  }, [refresh]);
  return <DataTable columns={columns} data={worlds}></DataTable>;
}

function ServerList({
  setAlertText,
}: {
  setAlertText: React.Dispatch<SetStateAction<string>>;
}) {
  const columns: ColumnDef<Server>[] = [
    {
      accessorKey: "id",
      header: "ID",
    },
    {
      accessorKey: "version",
      header: "Version",
    },
    {
      id: "actions",
      cell: ({ row }) => {
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuItem
                className={"red"}
                onClick={() => {
                  deleteServer(row.original.id);
                }}
              >
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];
  const [servers, setServers] = useState<Server[]>([]);
  const [refresh, setRefresh] = useState(1);

  function deleteServer(id: number) {
    del(`/servers/delete`, { id: id }).then((res) => {
      setAlertText(res);
      setRefresh((prev) => prev + 1);
    });
  }

  useEffect(() => {
    setAlertText("");
    get("/servers/list").then((res) => {
      setServers(res);
    });
  }, [refresh]);
  return <DataTable columns={columns} data={servers}></DataTable>;
}
