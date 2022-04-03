package main

import (
	"fmt"
	"strings"
	"time"
	"math"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/go-gl/gl/v4.6-core/gl"
	"math/rand"
)


var hexSize float32
var VBO uint32
var VAO uint32
var EBO uint32

func main() {
	tau := math.Pi * 2

	var windowWidth int32 = 700
	var windowHeight int32 = 700
	sin60 := float32(math.Sin(tau/6))

	xhexConsts := [6]float32{1,float32(math.Cos(tau/6)),float32(math.Cos(tau/3)),-1,float32(math.Cos(2*tau/3)),float32(math.Cos(5*tau/6))}
	yhexConsts := [6]float32{0,float32(math.Sin(tau/6)),float32(math.Sin(tau/3)),0, float32(math.Sin(2*tau/3)),float32(math.Sin(5*tau/6))}
	var hexConsts [12]float32
	for i:=0;i<6;i++ {
		hexConsts[i*2] = xhexConsts[i]
		hexConsts[i*2+1] = yhexConsts[i]
	}
	
	var heatTransferFactor float32 = 0.15
	
	hexSize = 1
	
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {panic(err)}

	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", 200,200,windowWidth,windowHeight,sdl.WINDOW_OPENGL)

	window.GLCreateContext()
	defer window.Destroy()

	gl.Init()
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println(version)

	
	vertexShaderSource := `#version 460 core
	layout (location = 0) in vec2 pos;
	

	void main() {	
		gl_Position = vec4(pos, 0.0, 1.0);
	}`+"\x00"

	fragShaderSource := `#version 460 core
	out vec4 fragColor;
	uniform vec4 color;

	void main() {
		fragColor = color;
	}`+"\x00"

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)

	csrc,free := gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, csrc, nil)
	free()
	gl.CompileShader(vertexShader)
	var status int32
	gl.GetShaderiv(vertexShader,gl.COMPILE_STATUS,&status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertexShader,gl.INFO_LOG_LENGTH,&logLength)
		log:=strings.Repeat("\x00",int(logLength+1))
		gl.GetShaderInfoLog(vertexShader,logLength,nil,gl.Str(log))
		panic("vertexshader\n"+log)
	}

	
	fragShader := gl.CreateShader(gl.FRAGMENT_SHADER)

	csrc2,free := gl.Strs(fragShaderSource)
	
	gl.ShaderSource(fragShader, 1, csrc2, nil)
	free()
	gl.CompileShader(fragShader)
	var status2 int32
	gl.GetShaderiv(fragShader,gl.COMPILE_STATUS,&status2)
	if status2 == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragShader,gl.INFO_LOG_LENGTH,&logLength)
		log:=strings.Repeat("\x00",int(logLength+1))
		gl.GetShaderInfoLog(fragShader,logLength,nil,gl.Str(log))
		panic("fragshader\n" + log)
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragShader)
	gl.LinkProgram(shaderProgram)

	var status3 int32
	gl.GetProgramiv(shaderProgram,gl.LINK_STATUS,&status3)
	if status3 == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram,gl.INFO_LOG_LENGTH,&logLength)
		log:=strings.Repeat("\x00",int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram,logLength,nil,gl.Str(log))
		panic("shader3\n" + log)
	}




	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragShader)

	gl.GenBuffers(1,&VBO)
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &EBO)

	updateHexRadius(15/float32(windowHeight)*2, &hexConsts)

	var hexMap [15][15][3]float32
	var hexMap2 [15][15][3]float32

	for i:=0;i<15;i++ {
		
		for ii:=0;ii<15;ii++ {
			a:=rand.Float32()
			b:=rand.Float32()
			hexMap[i][ii] = [3]float32{rand.Float32(),a*a,b*b*b}
		}
	}
	
	
	
	//var camx float32 = 0
	//var camy float32 = 0




	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		
		hexMap[0][14][0] = 1
		hexMap[14][0][0] = 0

		for i:=0;i<15;i++ {
			for ii:=0;ii<15;ii++ {
				for iii:=0;iii<3;iii++ {
					hexMap2[i][ii][iii] = hexMap[i][ii][iii]
					if i > 0 {hexMap2[i][ii][iii] -= (hexMap[i][ii][iii] - hexMap[i-1][ii][iii]) * shitpow(heatTransferFactor,iii)}
					if i < 14 {hexMap2[i][ii][iii] -= (hexMap[i][ii][iii] - hexMap[i+1][ii][iii]) * shitpow(heatTransferFactor,iii)}
					if ii > 0 {hexMap2[i][ii][iii] -= (hexMap[i][ii][iii] - hexMap[i][ii-1][iii]) * shitpow(heatTransferFactor,iii)}
					if ii < 14 {hexMap2[i][ii][iii] -= (hexMap[i][ii][iii] - hexMap[i][ii+1][iii]) * shitpow(heatTransferFactor,iii)}
					if i > 0 && ii < 14 {hexMap2[i][ii][iii] -= (hexMap[i][ii][iii] - hexMap[i-1][ii+1][iii]) * shitpow(heatTransferFactor,iii)}
					if i < 14 && ii > 0 {hexMap2[i][ii][iii] -= (hexMap[i][ii][iii] - hexMap[i+1][ii-1][iii]) * shitpow(heatTransferFactor,iii)}
				}
			} 
		}

		//render
		gl.ClearColor(0,0,0,0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		
		a := hexSize * 1.5
		b := hexSize*sin60
		c := hexSize*sin60*2
		var f float32 = 0
		var d float32 = 0
		for i:= 0; i < 15; i++ {
			var e float32 = 0
			for ii:= 0; ii < 15; ii++ {
				for iii:=0;iii<3;iii++ {
					hexMap[i][ii][iii] = hexMap2[i][ii][iii]
				}
				drawHex(f-0.5,e-d-0.5,hexMap[i][ii][0],hexMap[i][ii][1],hexMap[i][ii][2],0,&hexConsts, shaderProgram)
				
				e += c
			}
			d+= b
			f+=a
		}
		window.GLSwap()

		//100ms
		time.Sleep(100000000)
		//vertices[0] += 0.1
		//fmt.Println(vertices[0]);
	}

}

func updateHexRadius(newSize float32, b *[12]float32) {
	a := newSize/hexSize
	hexSize = newSize
	for i := 0; i < 12; i++ {
		(*b)[i] *= a
	}
}

//updateHexRadius(10)

func drawHex(x float32,y float32, red float32,green float32,blue float32, alpha float32, a *[12]float32, b uint32) {
	var vertices [12]float32
	for i:= 0; i < 12; i+=2 {
		vertices[i] = x + (*a)[i]
		vertices[i+1] = y + (*a)[i+1]
	}

	indices := []uint32{
		0,1,2,
		0,2,3,
		0,3,4,
		0,4,5}
	
	

	
	
	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4,gl.Ptr(&vertices[0]),gl.STATIC_DRAW)
	
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(&indices[0]), gl.STATIC_DRAW)
	
	gl.VertexAttribPointer(0,2,gl.FLOAT,false,2*4, nil)
	gl.EnableVertexAttribArray(0)
	
	tmp := "color\x00"

	vertexColorLocation := gl.GetUniformLocation(b, gl.Str(tmp));

	gl.UseProgram(b)
	gl.Uniform4f(vertexColorLocation, red, green, blue, alpha);

	gl.BindVertexArray(VAO)
	//gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.DrawElements(gl.TRIANGLES, 12, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}

func shitpow(a float32,b int)float32 {
	for b > 0 {
		b -- 
		a *= a
	}

	return a
}