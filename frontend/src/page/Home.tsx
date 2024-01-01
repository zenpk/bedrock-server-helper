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
import { useEffect, useState } from "react";
import { Card } from "@/components/Card.tsx";
import { del, get } from "@/util/request.ts";
import { Server } from "@/util/types.ts";
import { DataTable } from "@/components/DataTable.tsx";

export function Home() {
  const [worlds, setWorlds] = useState([]);
  return (
    <Card>
      <Tabs defaultValue="account" className="w-[400px]">
        <TabsList
          defaultValue={"worlds"}
          className={"flex justify-items-center"}
        >
          <TabsTrigger value="worlds">Worlds</TabsTrigger>
          <TabsTrigger value="servers">Servers</TabsTrigger>
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
      console.log(res);
      setRefresh((prev) => prev + 1);
    });
  }

  useEffect(() => {
    get("/servers/list").then((res) => {
      setServers(res);
    });
  }, [refresh]);
  return <DataTable columns={columns} data={servers}></DataTable>;
}
