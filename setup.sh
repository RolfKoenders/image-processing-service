
echo "1) Creating missing folders"
mkdir -p uploads/processed

echo "2) Installing node dependencies with npm"
cd api && npm i

echo "3) Ready! \nRun 'docker-compose up --build' to start everything up"