import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "custom_widget.dart";
import "sign_up.dart";

class LoginPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {

    return Scaffold(
      body: Center(
        child: SingleChildScrollView(
          child: Column(
            children: [
              Padding(padding: EdgeInsets.only(top: 30)),
              Text("AREA", style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold),),

              Text("Log in", style: TextStyle(fontSize: 40, fontWeight: FontWeight.bold),),

              CustomClassicForm(hint:'Email'),
              CustomPasswdForm(hint:'Password'),
              Text("Forgot your password?", style: TextStyle(fontSize: 15, fontWeight: FontWeight.bold),),

              SizedBox(height: 20),

              CustomButton(
                backgroundColor:Color.fromARGB(255, 34, 34, 34),
                text: AutoSizeText("Log in", style: TextStyle(fontSize: 30, fontWeight: FontWeight.bold, color: Colors.white),),
                onPressed: () {},
                size: Size(double.infinity, 60),
                padding: EdgeInsets.symmetric(horizontal: 5.0),
              ),

              Padding(padding: EdgeInsets.only(top: 15)),
              MyDivider(
                text: AutoSizeText("or", style: TextStyle(fontSize: 16, color: Colors.grey), maxLines: 1,),
                textPadding: EdgeInsets.symmetric(horizontal: 8.0),
              ),
              Padding(padding: EdgeInsets.only(top: 15)),

              CustomButton(
                backgroundColor:Color.fromARGB(255, 232, 74, 74),
                text: AutoSizeText("Continue with google", style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold, color: Colors.white), maxLines: 1,),
                onPressed: () {},
                size: Size(double.infinity, 60),
                padding: EdgeInsets.symmetric(horizontal: 5.0),
                textPaddingLeft: 30,
                image: Image.asset('assets/icon/google_logo.webp', width: 24, height: 24),
              ),

              Padding(padding: EdgeInsets.only(top: 10)),
              CustomButton(
                backgroundColor:Color.fromARGB(255, 59, 87, 157),
                text: AutoSizeText("Continue with Facebook", style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold, color: Colors.white), maxLines: 1,),
                onPressed: () {},
                size: Size(double.infinity, 60),
                padding: EdgeInsets.symmetric(horizontal: 5.0),
                textPaddingLeft: 30,
                image: Image.asset('assets/icon/facebook_icon.png', width: 24, height: 24),
              ),

              Padding(padding: EdgeInsets.only(top: 10)),
              CustomButton(
                backgroundColor:Color.fromARGB(255, 169, 169, 169),
                text: AutoSizeText("Continue with X", style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold, color: Colors.white), maxLines: 1,),
                onPressed: () {},
                size: Size(double.infinity, 60),
                padding: EdgeInsets.symmetric(horizontal: 5.0),
                textPaddingLeft: 30,
                image: Image.asset('assets/icon/X-Logo.png', width: 40, height: 40),
              ),

              Padding(padding: EdgeInsets.only(top: 10)),
              CustomButton(
                backgroundColor:Color.fromARGB(255, 34, 34, 34),
                text: AutoSizeText("Continue with Apple", style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold, color: Colors.white), maxLines: 1,),
                onPressed: () {},
                size: Size(double.infinity, 60),
                padding: EdgeInsets.symmetric(horizontal: 5.0),
                textPaddingLeft: 30,
                image: Image.asset('assets/icon/apple_logo.png', width: 29, height: 29),
              ),

              Padding(padding: EdgeInsets.only(top: 30)),

              ElevatedButton(
                child: const Text("New to AREA? Sign up here. ", style: TextStyle(fontSize: 15, fontWeight: FontWeight.bold, color: Colors.black),),
                onPressed: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(builder: (context) => SignUp()),
                  );
                },
              ),

            ],
          ),
        ),
      ),
    );
  }
}