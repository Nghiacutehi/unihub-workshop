import { Navbar } from "@/components/student/navbar"
import { WorkshopFilters } from "@/components/student/workshop-filters"
import { WorkshopGrid } from "@/components/student/workshop-grid"

export default function HomePage() {
  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-2xl font-semibold text-foreground sm:text-3xl">
            Khám phá Workshop
          </h1>
          <p className="mt-2 text-muted-foreground">
            Tìm và đăng ký các workshop phù hợp với bạn
          </p>
        </div>

        {/* Filters Section */}
        <div className="mb-8">
          <WorkshopFilters />
        </div>

        {/* Workshop Grid */}
        <WorkshopGrid />
      </main>
    </div>
  )
}
