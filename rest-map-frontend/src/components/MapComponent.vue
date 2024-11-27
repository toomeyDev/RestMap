<!-- src/components/Map.vue -->
<template>
  <div>
    <div id="controls">
      <input v-model="chainName" placeholder="Enter restaurant chain" />
      <button @click="searchRestaurants">Search</button>     
    </div> 
    <div id="map-component" style="width: 100%; height: calc(100vh - 50px);"></div>
  </div>
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
        chainName: '',
        markers: [],
      };
    },
    mounted() {
      mapboxgl.accessToken = process.env.VUE_APP_MAPBOX_ACCESS_TOKEN;

      // Request user's location
      this.getUserLocation()
        .then(coords => {
          this.userLocation = coords;
          this.initializeMap(coords);
        })
        .catch(error => {
          console.error('Error getting user location:', error);
          // Fallback to NYC as a default
          const fallbackCoords = { latitude: 40.7128, longitude: -74.0060 };
          this.initializeMap(fallbackCoords);
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
            zoom: 12,
          });

          // Add user's location marker
          new mapboxgl.Marker({ color: 'blue' })
            .setLngLat([coords.longitude, coords.latitude])
            .setPopup(new mapboxgl.Popup().setText('Your Location'))
            .addTo(this.map);
        },
        searchRestaurants() {
          if (!this.userLocation) {
            console.error('User location is not available.');
            return;
          }
          if (!this.chainName) {
            console.error('Please enter a restaurant chain name.');
            return;
          }
          
          // Remove any existing markers
          this.markers.forEach(marker => marker.remove());
          this.markers = [];
          // Fetch new restaurants w/ the current query          
          this.fetchRestaurants(this.userLocation);
        },
        fetchRestaurants(coords) {
          axios
            .get('http://localhost:8080/restaurants', {
              params: {
                latitude: coords.latitude,
                longitude: coords.longitude,
                chain: this.chainName
              },
            })
            .then(response => {
              const restaurants = response.data.restaurants;
              restaurants.forEach(restaurant => {
                const popupContent = `
                  <strong>${restaurant.name}</strong><br />
                  ${restaurant.address || 'Address not available'}<br/>
                  ${restaurant.category || 'Category not available'}
                `;
                const marker = new mapboxgl.Marker()
                  .setLngLat([restaurant.longitude, restaurant.latitude])
                  .setPopup(new mapboxgl.Popup().setHTML(popupContent))
                  .addTo(this.map);

                this.markers.push(marker);
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
  #controls {
    position: absolute;
    z-index: 1;
    background: white;
    padding: 10px;
    width: 100%;
  }
  #map-component {
    position: relative;
  }
</style>
