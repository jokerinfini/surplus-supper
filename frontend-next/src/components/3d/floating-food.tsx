'use client'

import { Canvas, useFrame } from '@react-three/fiber'
import { useRef } from 'react'
import { OrbitControls } from '@react-three/drei'
import * as THREE from 'three'

// Simple rotating cube component
function RotatingCube() {
  const meshRef = useRef<THREE.Mesh>(null)

  useFrame((state) => {
    if (meshRef.current) {
      meshRef.current.rotation.x = Math.sin(state.clock.elapsedTime) * 0.3
      meshRef.current.rotation.y += 0.01
    }
  })

  return (
    <mesh ref={meshRef} position={[0, 0, 0]}>
      <boxGeometry args={[1, 1, 1]} />
      <meshStandardMaterial color="#22c55e" />
    </mesh>
  )
}

// Food items
function FoodItems() {
  const pizzaRef = useRef<THREE.Mesh>(null)
  const burgerRef = useRef<THREE.Mesh>(null)
  const sushiRef = useRef<THREE.Mesh>(null)

  useFrame((state) => {
    if (pizzaRef.current) {
      pizzaRef.current.rotation.y += 0.01
      pizzaRef.current.position.y = Math.sin(state.clock.elapsedTime * 0.5) * 0.2
    }
    if (burgerRef.current) {
      burgerRef.current.rotation.y += 0.01
      burgerRef.current.position.y = Math.sin(state.clock.elapsedTime * 0.5 + 1) * 0.2
    }
    if (sushiRef.current) {
      sushiRef.current.rotation.y += 0.01
      sushiRef.current.position.y = Math.sin(state.clock.elapsedTime * 0.5 + 2) * 0.2
    }
  })

  return (
    <>
      {/* Pizza */}
      <mesh ref={pizzaRef} position={[-2, 0, 0]}>
        <cylinderGeometry args={[0.5, 0.5, 0.1, 32]} />
        <meshStandardMaterial color="#f97316" />
      </mesh>
      
      {/* Burger */}
      <mesh ref={burgerRef} position={[2, 0, 0]}>
        <boxGeometry args={[0.6, 0.4, 0.6]} />
        <meshStandardMaterial color="#dc2626" />
      </mesh>
      
      {/* Sushi */}
      <mesh ref={sushiRef} position={[0, 2, 0]}>
        <cylinderGeometry args={[0.2, 0.2, 0.1, 8]} />
        <meshStandardMaterial color="#059669" />
      </mesh>
    </>
  )
}

// Main scene
function Scene() {
  return (
    <>
      <ambientLight intensity={0.5} />
      <pointLight position={[10, 10, 10]} intensity={1} />
      <pointLight position={[-10, -10, -10]} intensity={0.5} />
      
      <FoodItems />
    </>
  )
}

export function FloatingFoodScene() {
  return (
    <div className="w-full h-96 relative border-2 border-green-200 rounded-lg">
      <Canvas
        camera={{ position: [0, 0, 5], fov: 75 }}
        style={{ background: 'linear-gradient(135deg, #f0fdf4 0%, #dcfce7 50%, #bbf7d0 100%)' }}
      >
        <Scene />
        <OrbitControls 
          enableZoom={true}
          enablePan={true}
          autoRotate
          autoRotateSpeed={0.5}
        />
      </Canvas>
      
      {/* Fallback text if 3D doesn't load */}
      <div className="absolute inset-0 flex items-center justify-center text-green-600 font-semibold">
        <div className="text-center">
          <div className="text-4xl mb-2">🍽️</div>
          <div>3D Food Scene</div>
        </div>
      </div>
    </div>
  )
}
