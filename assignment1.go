package main

import (
"fmt"
"gobot.io/x/gobot"
"gobot.io/x/gobot/drivers/aio"
g "gobot.io/x/gobot/platforms/dexter/gopigo3"
"gobot.io/x/gobot/platforms/raspi"
"os"
"time"
)

//robotRunLoop is the main function for the robot, the gobot framework
//will spawn a new thread in the NewRobot factory function and run this
//function in that new thread. Do all of your work in this function and
//in other functions that this function calls. don't read from sensors or
//use actuators from main or you will get a panic.
func robotRunLoop(lightSensorLeft *aio.GroveLightSensorDriver, lightSensorRight *aio.GroveLightSensorDriver, gopigo3 *g.Driver) {

	var sensorRange int = 200
	var minLight int = 800

	for {
		sensorValLeft, err := lightSensorLeft.Read()

		if err != nil {
			fmt.Errorf("Error reading light sensor %+v", err)
		}
		sensorValRight, err := lightSensorRight.Read()
		if err != nil {
			fmt.Errorf("Error reading light sensor %+v", err)
		}

		if sensorValLeft >= 3000 && sensorValRight >= 3000 {
			gopigo3.Halt()
			os.Exit(0)
		}
		// if facing wrong way -- rotate
		// if both sensors are within tolerance BUT light is below 1500
		if (sensorValLeft < 1200 && sensorValLeft > minLight) && (sensorValRight < 1200 && sensorValRight > minLight) {
			gopigo3.SetMotorDps(g.MOTOR_LEFT, -50)
			gopigo3.SetMotorDps(g.MOTOR_RIGHT, 50)
			listValues(lightSensorLeft, lightSensorRight)
			fmt.Print("Sensors below 1500")

		} else if sensorValLeft < 3000 && sensorValRight < 3000 {
			listValues(lightSensorLeft, lightSensorRight)

			if (sensorValLeft - sensorValRight) > sensorRange {
				gopigo3.SetMotorDps(g.MOTOR_LEFT, -50)
				gopigo3.SetMotorDps(g.MOTOR_RIGHT, 50)
				fmt.Printf("Turning Left...")

			} else if (sensorValRight - sensorValLeft) > sensorRange {
				gopigo3.SetMotorDps(g.MOTOR_LEFT, 50)
				gopigo3.SetMotorDps(g.MOTOR_RIGHT, -50)
				fmt.Printf("Turning Right...")

			} else if (sensorValRight-sensorValLeft) <= sensorRange || (sensorValLeft-sensorValRight) <= sensorRange {
				gopigo3.SetMotorDps(g.MOTOR_RIGHT+g.MOTOR_LEFT, 100)
				fmt.Printf("Proceeding Forward...")

			}

		} else {
			fmt.Printf("You got problems, kid")

		}
		fmt.Println("Light Value is ", sensorValLeft)
		fmt.Println("Light Value is ", sensorValRight)
		time.Sleep(time.Second)
	}
}
func listValues(lightSensorLeft *aio.GroveLightSensorDriver, lightSensorRight *aio.GroveLightSensorDriver){
	fmt.Print(lightSensorLeft)
	fmt.Print(lightSensorRight)
	fmt.Print("\n")
}

func main() {
	//We create the adaptors to connect the GoPiGo3 board with the Raspberry Pi 3
	//also create any sensor drivers here
	raspiAdaptor := raspi.NewAdaptor()
	gopigo3 := g.NewDriver(raspiAdaptor)
	lightSensorLeft := aio.NewGroveLightSensorDriver(gopigo3, "AD_2_1") //AnalogDigital Port 1 is "AD_1_1" this is port 2
	lightSensorRight := aio.NewGroveLightSensorDriver(gopigo3, "AD_1_1")
	//end create hardware drivers

	//here we create an anonymous function assigned to a local variable
	//the robot framework will create a new thread and run this function
	//I'm calling my robot main loop here. Pass any of the variables we created
	//above to that function if you need them
	mainRobotFunc := func() {
		robotRunLoop(lightSensorLeft, lightSensorRight, gopigo3)
	}

	//this is the crux of the gobot framework. The factory function to create a new robot
	//struct (go uses structs and not objects) It takes four parameters
	robot := gobot.NewRobot("gopigo3sensorChecker", //first a name
		[]gobot.Connection{raspiAdaptor},                  //next a slice of connections to one or more robot controllers
		[]gobot.Device{gopigo3, lightSensorLeft, lightSensorRight}, //next a slice of one or more sensors and actuators for the robots
		mainRobotFunc, //the variable holding the function to run in a new thread as the main function
	)

	robot.Start() //actually run the function
}


