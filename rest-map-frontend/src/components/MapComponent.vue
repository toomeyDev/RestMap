<!-- src/components/Map.vue -->
<template>
  <div id="map-component" style="width: 100%; height: 100vh"></div>
</template>

<script>
  import mapboxgl from "mapbox-gl";
  import axios from 'axios';

  export default {
    name: 'Map-Component',
    data () {
      return {
        map: null,
        userLocation: null,
      };
    },
    mounted() {
      mapboxgl.accessToken = 'pk.eyJ1IjoibXRvb21leS1kZXYiLCJhIjoiY200MDRiOWxqMTlzODJ3cHlrejE3NjgzNyJ9.AqCgaLym-Ym9jBSrsj6aGg';
      
      // Request user's location
      this.getUserLocation()
        .then(coords => {
          this.userLocation = coords;
          this.initializeMap(coords);
          this.fetchRestaurants(coords);
        })
        .catch(error => {
          console.error('Error getting user location:', error);
          // Fallback to NYC as a default
          const fallbackCoords = { latitude: 40.7128, longitude: -74.0060 };
          this.initializeMap(fallbackCoords);
          this.fetchRestaurants(fallbackCoords);
        });
      },
      methods: {
        getUserLocation() {
          return new Promise((resolve, reject) => {
            if (navigator.geolocation) {
              navigator.geolocation.getCurrentPosition(
                position => {
                  const { latitude, longitude } = position.coords;
                  resolve({ latitude, longitude });
                },
                error => {
                  reject(error);
                }
              );
            } else {
              reject(new Error('Geolocation is not supported by this browser.'));
            }
          });
        },
        initializeMap(coords) {
          this.map = new mapboxgl.Map({
            container: 'map-component',
            style: 'mapbox://styles/mapbox/streets-v11',
            center: [coords.longitude, coords.latitude],
            zoom: 12
          });

          // Add user's location marker
          new mapboxgl.Marker({ color: 'blue' })
            .setLngLat([coords.longitude, coords.latitude])
            .setPopup(new mapboxgl.Popup().setText('Your Location'))
            .addTo(this.map);
        },
        fetchRestaurants(coords) {
          const chainName = 'Red Lobster'; // use as an example for now
          axios
            .get('http://localhost:8080/restaurants', {
              params: {
                latitude: coords.latitude,
                longitude: coords.longitude,
                chain: chainName
              },
            })
            .then(response => {
              const restaurants = response.data.restaurants;
              restaurants.forEach(restaurant => {
                new mapboxgl.Marker()
                  .setLngLat([restaurant.longitude, restaurant.latitude])
                  .setPopup(new mapboxgl.Popup().setText(restaurant.name))
                  .addTo(this.map);
              });
            })
            .catch(error => {
              console.error('Error fetching restaurant data: ', error);
            });
        },
      },
    };
</script>

<style>
  #map-component {
    position: relative;
  }
</style>
