import { GreenMo } from '../lib/query/GreenMo';
import axios from 'axios';
jest.mock('axios');

test('the location of chargable cars is fetched', async () => {
    const car1 = {
        id: 1,
        position: {
            coordinates: [2.123456, 1.123456]
        },
        stateOfCharge: 30,
    };
    const car2 = {
        id: 2,
        position: {
            coordinates: [4.123456, 3.123456]
        },
        stateOfCharge: 50,
    };
    const data = [car1, car2];

    (axios.get as jest.Mock).mockImplementation(() =>
        Promise.resolve({ status: 200, data: data })
    );

    const pos1 = { lat: 1.123456, lon: 2.123456 };
    const pos2 = { lat: 3.123456, lon: 4.123456 };
    const params = {
        lon1: `${pos1.lon}`,
        lat1: `${pos1.lat}`,
        lon2: `${pos2.lon}`,
        lat2: `${pos2.lat}`,
    };
    const greenMo = new GreenMo(40);
    const cars = await greenMo.query(params);
    expect(cars).toEqual([{ lat: car1.position.coordinates[1], lon: car1.position.coordinates[0] }]);
});
