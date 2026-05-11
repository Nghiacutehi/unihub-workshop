"use client"

import { Calendar, MapPin, User } from "lucide-react"
import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Progress } from "@/components/ui/progress"
import { Button } from "@/components/ui/button"

export interface Workshop {
  id: string
  title: string
  speaker: string
  speakerTitle?: string
  date: string
  time: string
  location: string
  capacity: number
  availableSeats: number
  ticketType: "free" | "paid"
  price?: number
  category: string
  imageUrl?: string
}

interface WorkshopCardProps {
  workshop: Workshop
}

export function WorkshopCard({ workshop }: WorkshopCardProps) {
  const filledPercentage = Math.round(
    ((workshop.capacity - workshop.availableSeats) / workshop.capacity) * 100
  )
  const isAlmostFull = filledPercentage >= 80
  const isFull = workshop.availableSeats === 0

  return (
    <Card className="group flex h-full flex-col overflow-hidden transition-shadow hover:shadow-md">
      {/* Image placeholder */}
      <div className="relative h-40 w-full bg-muted">
        {workshop.imageUrl ? (
          <img
            src={workshop.imageUrl}
            alt={workshop.title}
            className="h-full w-full object-cover"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-gradient-to-br from-primary/5 to-primary/10">
            <span className="text-4xl font-bold text-primary/20">
              {workshop.title.charAt(0)}
            </span>
          </div>
        )}
        {/* Category Badge */}
        <Badge
          variant="secondary"
          className="absolute left-3 top-3 bg-card/90 backdrop-blur-sm"
        >
          {workshop.category}
        </Badge>
        {/* Ticket Type Badge */}
        <Badge
          className={`absolute right-3 top-3 ${
            workshop.ticketType === "free"
              ? "bg-success text-success-foreground"
              : "bg-primary text-primary-foreground"
          }`}
        >
          {workshop.ticketType === "free"
            ? "Miễn phí"
            : `${workshop.price?.toLocaleString("vi-VN")}đ`}
        </Badge>
      </div>

      <CardHeader className="gap-2 pb-2">
        <h3 className="line-clamp-2 text-base font-semibold leading-tight text-foreground group-hover:text-primary transition-colors">
          {workshop.title}
        </h3>
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <User className="h-4 w-4 shrink-0" />
          <span className="truncate">{workshop.speaker}</span>
        </div>
      </CardHeader>

      <CardContent className="flex flex-col gap-3 pb-4">
        <div className="flex flex-col gap-1.5 text-sm text-muted-foreground">
          <div className="flex items-center gap-2">
            <Calendar className="h-4 w-4 shrink-0" />
            <span>
              {workshop.date} • {workshop.time}
            </span>
          </div>
          <div className="flex items-center gap-2">
            <MapPin className="h-4 w-4 shrink-0" />
            <span className="truncate">{workshop.location}</span>
          </div>
        </div>

        {/* Seat Progress */}
        <div className="space-y-1.5">
          <div className="flex items-center justify-between text-xs">
            <span className="text-muted-foreground">Chỗ ngồi</span>
            <span
              className={`font-medium ${
                isFull
                  ? "text-destructive"
                  : isAlmostFull
                  ? "text-warning"
                  : "text-muted-foreground"
              }`}
            >
              {workshop.availableSeats}/{workshop.capacity}
            </span>
          </div>
          <Progress
            value={filledPercentage}
            className={`h-1.5 ${
              isFull
                ? "[&>div]:bg-destructive"
                : isAlmostFull
                ? "[&>div]:bg-warning"
                : "[&>div]:bg-success"
            }`}
          />
        </div>
      </CardContent>

      <CardFooter className="mt-auto pt-0">
        <Button
          className="h-10 w-full font-medium"
          disabled={isFull}
        >
          {isFull ? "Hết chỗ" : "Đăng ký ngay"}
        </Button>
      </CardFooter>
    </Card>
  )
}
