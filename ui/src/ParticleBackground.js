import React, { useEffect, useMemo, useRef } from "react";
import Particles, { initParticlesEngine } from "@tsparticles/react";
import { loadSlim } from "@tsparticles/slim";

const ParticleBackground = React.memo(() => {
  const particlesLoadedRef = useRef(false);  // Ref to track if particles have been loaded

  useEffect(() => {
    // Initialize particles only once
    if (!particlesLoadedRef.current) {
      initParticlesEngine(async (engine) => {
        await loadSlim(engine);  // Initialize particles
      }).then(() => {
        particlesLoadedRef.current = true;  // Set particles as loaded
      });
    }
  }, []);

  const particlesLoaded = (container) => {
    console.log(container);
  };

  const options = useMemo(
    () => ({
      background: {
        color: {
          value: "#000000", // Background color
        },
      },
      fpsLimit: 120,
      particles: {
        color: {
          value: "#FFFFF2", // Particle color
        },
        links: {
          color: "#FFFFF2", // Link color
          distance: 200,
          enable: true,
          opacity: 0.1,
          width: 2,
        },
        move: {
          direction: "none",
          enable: true,
          speed: 1,
        },
        number: {
          value: 50, // Number of particles
        },
        opacity: {
          value: 0.1,
        },
        shape: {
          type: "circle",
        },
        size: {
          value: { min: 1, max: 5 },
        },
      },
      detectRetina: true,
    }),
    [],
  );

  return (
    <Particles
      id="tsparticles"
      particlesLoaded={particlesLoaded}
      options={options}
    />
  );
});

export default ParticleBackground;

