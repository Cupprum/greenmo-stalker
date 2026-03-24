import { PositionQuery, Position } from './PositionQuery';

interface Car extends Position {
    id: number;
    stateOfCharge: number;
    position: {
        coordinates: [number, number];
    };
}


export class GreenMo extends PositionQuery<Car> {
    protocol = 'https';
    hostname = 'platform.api.gourban.services';
    endpoint = 'v1/hb98ga69/front/vehicles'; // hb98ga69 should be tenant ID of greenmobility
    desiredFuelLevel: number;

    constructor(desiredFuelLevel: number) {
        super();
        this.desiredFuelLevel = desiredFuelLevel;
    }

    protected filter(objs: Car[]): Car[] {
        return objs.filter(
            (car: Car) => car.stateOfCharge <= this.desiredFuelLevel
        );
    }

    protected map(objs: Car[]): Position[] {
        return objs.map((car: Car) => {
            return {
                lat: car.position.coordinates[1],
                lon: car.position.coordinates[0],
            };
        });
    }
}
