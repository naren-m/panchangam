// swift-tools-version: 6.0

import PackageDescription

let package = Package(
    name: "PanchangamShared",
    platforms: [
        .iOS(.v17),
        .watchOS(.v10),
        .macOS(.v14)
    ],
    products: [
        .library(name: "PanchangamShared", targets: ["PanchangamShared"]),
        .executable(name: "PanchangamSharedChecks", targets: ["PanchangamSharedChecks"])
    ],
    targets: [
        .target(name: "PanchangamShared"),
        .executableTarget(
            name: "PanchangamSharedChecks",
            dependencies: ["PanchangamShared"],
            path: "Checks/PanchangamSharedChecks"
        )
    ]
)
