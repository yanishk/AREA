import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:convert';
import "package:http/src/response.dart";
import 'dart:developer' as developer;

import "custom_widget.dart";
import "sign_up.dart";
import "exploreAll.dart";
import "globalData.dart";
import "loadServicesList.dart";
import 'apiCall/postApi.dart';

class LoginPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: APICall(),
    );
  }
}

class APICall extends StatefulWidget {
  @override
  _APICallState createState() => _APICallState();
}

class _APICallState extends State<APICall> {
  final TextEditingController emailController = new TextEditingController(text: "test@test.com");
  final TextEditingController passwordController = new TextEditingController(text: "test");
  String result = '';
  String apiRoute = 'login';

  @override
  void dispose() {
    emailController.dispose();
    passwordController.dispose();
    super.dispose();
  }

  @override
  void initState() {
    super.initState();
    _checkJWToken();
  }

  void _checkJWToken() {
    if (Globaldata.JWToken.isNotEmpty) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(builder: (context) => ExplorePageAll()),
        );
      });
    }
  }

  Future<void> saveJWToken(Response response) async {
    String? cookiesBrute = response.headers["set-cookie"];
    String JWToken = "";
    if (cookiesBrute != null) {
      JWToken = cookiesBrute.split(';').first;
      setState(() {
        result = "";
        developer.log("JWToken : {${JWToken}}");
        Globaldata.JWToken = JWToken;
        Navigator.push(
          Globaldata.myContext,
          MaterialPageRoute(builder: (context) => LoadServicePage()),
        );
      });
    } else {
      setState(() {
        result = 'Failed to get the token';
      });
    }
  }

  Future<void> postLoginData() async {
    final response = await PostApi.post(apiRoute, jsonEncode(<String, dynamic>{'email': emailController.text,'password': passwordController.text,}),);
    if (response.statusCode == 0) {
      setState(() {
        result = response.body;
      });
      return;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "postLoginData FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      setState(() {
        result = '${responseData['error']}';
      });
      return;
    }
    saveJWToken(response);
  }

  @override
  Widget build(BuildContext context) {
    Globaldata.myContext = context;
    return Scaffold(
      body: Center(
        child: SingleChildScrollView(
          child: Column(
            children: [
              Padding(padding: EdgeInsets.only(top: 30)),
              Text(
                "AREA",
                style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold),
              ),
              Text(
                "Log in",
                style: TextStyle(fontSize: 40, fontWeight: FontWeight.bold),
              ),
              Text(
                result,
                style: TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: Colors.red),
              ),
              CustomClassicForm(
                hint: 'Email',
                controller: emailController,
              ),
              CustomClassicForm(
                  hint: 'Password',
                  controller: passwordController,
                  obscureText: true),
              Text(
                "Forgot your password?",
                style: TextStyle(fontSize: 15, fontWeight: FontWeight.bold),
              ),
              SizedBox(height: 20),
              CustomButton(
                backgroundColor: Color.fromARGB(255, 34, 34, 34),
                text: AutoSizeText(
                  "Log in",
                  style: TextStyle(
                      fontSize: 30,
                      fontWeight: FontWeight.bold,
                      color: Colors.white),
                ),
                onPressed: postLoginData,
                size: Size(double.infinity, 70),
                padding: EdgeInsets.symmetric(horizontal: 15.0),
              ),
              Padding(padding: EdgeInsets.only(top: 15)),
              MyDivider(
                text: AutoSizeText(
                  "or",
                  style: TextStyle(fontSize: 16, color: Colors.grey),
                  maxLines: 1,
                ),
                textPadding: EdgeInsets.symmetric(horizontal: 8.0),
              ),
              Padding(padding: EdgeInsets.only(top: 15)),
              CustomButton(
                backgroundColor: Color.fromARGB(255, 232, 74, 74),
                text: AutoSizeText(
                  "Continue with google",
                  style: TextStyle(
                      fontSize: 25,
                      fontWeight: FontWeight.bold,
                      color: Colors.white),
                  maxLines: 1,
                ),
                onPressed: () {},
                size: Size(double.infinity, 60),
                padding: EdgeInsets.symmetric(horizontal: 15.0),
                textPaddingLeft: 30,
                image: Image.asset('assets/icon/google_logo.webp',
                    width: 24, height: 24),
              ),
              Padding(padding: EdgeInsets.only(top: 10)),
              CustomButton(
                backgroundColor: Color.fromARGB(255, 59, 87, 157),
                text: AutoSizeText(
                  "Continue with Facebook",
                  style: TextStyle(
                      fontSize: 25,
                      fontWeight: FontWeight.bold,
                      color: Colors.white),
                  maxLines: 1,
                ),
                onPressed: () {},
                size: Size(double.infinity, 60),
                padding: EdgeInsets.symmetric(horizontal: 15.0),
                textPaddingLeft: 30,
                image: Image.asset('assets/icon/facebook_icon.png',
                    width: 24, height: 24),
              ),
              Padding(padding: EdgeInsets.only(top: 10)),
              CustomButton(
                backgroundColor: Color.fromARGB(255, 169, 169, 169),
                text: AutoSizeText(
                  "Continue with X",
                  style: TextStyle(
                      fontSize: 25,
                      fontWeight: FontWeight.bold,
                      color: Colors.white),
                  maxLines: 1,
                ),
                onPressed: () {},
                size: Size(double.infinity, 60),
                padding: EdgeInsets.symmetric(horizontal: 15.0),
                textPaddingLeft: 30,
                image: Image.asset('assets/icon/X-Logo.png',
                    width: 40, height: 40),
              ),
              Padding(padding: EdgeInsets.only(top: 10)),
              CustomButton(
                backgroundColor: Color.fromARGB(255, 34, 34, 34),
                text: AutoSizeText(
                  "Continue with Apple",
                  style: TextStyle(
                      fontSize: 25,
                      fontWeight: FontWeight.bold,
                      color: Colors.white),
                  maxLines: 1,
                ),
                onPressed: () {
                  Navigator.push(
                    Globaldata.myContext,
                    MaterialPageRoute(builder: (context) => LoadServicePage()),
                  );
                },
                size: Size(double.infinity, 60),
                padding: EdgeInsets.symmetric(horizontal: 15.0),
                textPaddingLeft: 30,
                image: Image.asset('assets/icon/apple_logo.png',
                    width: 29, height: 29),
              ),
              Padding(padding: EdgeInsets.only(top: 30)),
              CustomTextButton(
                text: Text(
                  "New to AREA? Sign up here.",
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.bold,
                    color: Colors.black,
                  ),
                ),
                onPressed: () {
                  Navigator.push(
                    Globaldata.myContext,
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
